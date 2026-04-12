import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AgGridAngular } from 'ag-grid-angular';
import { ColDef } from 'ag-grid-community';
import { ModuleRegistry, AllCommunityModule } from 'ag-grid-community';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { WorkerService, WorkerShift } from '../services/worker.service';
import { Observable, tap } from 'rxjs';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';

ModuleRegistry.registerModules([AllCommunityModule]);

interface Shift {
  date: string;
  startTime: string;
  endTime: string;
  status: 'Assigned' | 'Completed';
  hourlyWage: number;
}

@Component({
  selector: 'app-worker-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    AgGridAngular,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatChipsModule
  ],
  templateUrl: './worker-dashboard.component.html',
  styleUrls: ['./worker-dashboard.component.scss']
})
export class WorkerDashboardComponent implements OnInit {

  // Dynamic User Data (Replaced Hardcoded Values)
  currentWorkerId: number = 0;
  currentCompanyId: number = 0;
  currentWorkerName: string = '';

  currentDateTime = '';
  isProfileMenuOpen = false;

  myShifts$!: Observable<WorkerShift[]>;
  nextShift: WorkerShift | null = null;
  assignedShifts: WorkerShift[] = [];

  totalHours = 0;
  estEarnings = 0;
  hourlyWage = 15;
  weeklyTarget = 40;
  dailyTarget = 8;
  weeklyTotal = 0;
  weeklyProgress = 0;
  weekRangeLabel = '';
  weeklyDays: Array<{ label: string; date: Date; hours: number; percent: number }> = [];
  weekOffset = 0;
  private latestShifts: WorkerShift[] = [];
  isReleasing: Record<number, boolean> = {};
  releaseError = '';

  constructor(
    private workerService: WorkerService,
    private authService: AuthService,
    private router: Router
  ) { }

  colDefs: ColDef[] = [
    // 1. Extract the Date from the start_time ISO string
    {
      headerName: 'Date',
      valueGetter: params => params.data.start_time,
      valueFormatter: params => new Date(params.value).toLocaleDateString(),
      flex: 1
    },
    // 2. Map to snake_case 'start_time'
    {
      field: 'start_time',
      headerName: 'Start',
      valueFormatter: params => new Date(params.value).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      flex: 1
    },
    // 3. Map to snake_case 'end_time'
    {
      field: 'end_time',
      headerName: 'End',
      valueFormatter: params => new Date(params.value).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      flex: 1
    },
    // 4. Duration of the shift
    {
      headerName: 'Duration',
      valueGetter: (params) => {
        const start = new Date(params.data.start_time).getTime();
        const end = new Date(params.data.end_time).getTime();
        const hours = (end - start) / (1000 * 60 * 60);
        return hours.toFixed(1) + ' hrs';
      },
      flex: 1
    },
    // 5. Update status check for lowercase 'assigned'
    {
      field: 'status',
      flex: 1,
      cellStyle: (params) => {
        // Go backend returns lowercase "assigned"
        const isAssigned = params.value?.toLowerCase() === 'assigned';
        return {
          color: isAssigned ? '#0f5f5c' : '#666',
          fontWeight: 'bold'
        };
      }
    },
    // 6. Wage is not currently in the Shift API response, 
    // so we can use your hardcoded dashboard value for now.
    {
      headerName: 'Wage ($/hr)',
      valueGetter: () => this.hourlyWage,
      flex: 1
    },
  ];

  ngOnInit() {
    this.updateDateTime();

    // 1. Ask the single source of truth for the user
    const user = this.authService.currentUserValue;

    // 2. Safely assign variables and fetch data
    if (user && user.id) {
      this.currentWorkerId = user.id;
      this.currentWorkerName = user.name || 'Worker';
      this.currentCompanyId = user.company_id || 1;

      console.log(`Worker Dashboard initialized for User ID: ${this.currentWorkerId}`);
      this.weekOffset = 0;

      // 3. Load shifts using the REAL dynamic ID
      this.reloadShifts();
    } else {
      console.error('CRITICAL: No valid user state found. Cannot load worker shifts.');
    }
  }

  private reloadShifts() {
    const { start, end } = this.getUpcomingWindow();
    this.myShifts$ = this.workerService.getMyShifts(this.currentWorkerId, start, end).pipe(
      tap((data) => {
        const safeData = data || [];
        this.latestShifts = safeData;
        console.log("Stream received data:", safeData);
        this.calculateStats(safeData);
        this.findNextShift(safeData);
        this.assignedShifts = this.getAssignedShifts(safeData);
        this.calculateWeeklyAssigned(safeData);
      })
    );
  }

  private getUpcomingWindow() {
    const start = new Date();
    start.setHours(0, 0, 0, 0);
    const end = new Date(start);
    end.setDate(end.getDate() + 28);
    return {
      start: start.toISOString(),
      end: end.toISOString()
    };
  }

  releaseShift(shift: WorkerShift) {
    const shiftId = Number(shift?.shift_id ?? shift?.id);
    if (!shiftId || this.isReleasing[shiftId]) return;

    this.releaseError = '';
    this.isReleasing[shiftId] = true;

    this.workerService.releaseShift(shiftId).subscribe({
      next: () => {
        this.isReleasing[shiftId] = false;
        this.reloadShifts();
      },
      error: (err) => {
        console.error('Failed to release shift', err);
        this.isReleasing[shiftId] = false;
        this.releaseError = 'Failed to release shift. Please try again.';
      }
    });
  }

  calculateStats(shifts: WorkerShift[]) {
    this.totalHours = 0;

    shifts.forEach(shift => {
      // Safety check for missing times
      if (!shift.start_time || !shift.end_time) return;

      const start = new Date(shift.start_time).getTime();
      const end = new Date(shift.end_time).getTime();
      const durationHours = (end - start) / (1000 * 60 * 60);

      if (durationHours > 0) {
        this.totalHours += durationHours;
      }
    });

    this.estEarnings = this.totalHours * this.hourlyWage;
  }

  findNextShift(shifts: WorkerShift[]) {
    const now = new Date().getTime();

    // safety check
    if (!shifts) return;

    const assignedShifts = this.getAssignedShifts(shifts);

    // Filter for future assigned shifts and sort by start time
    const futureShifts = assignedShifts
      .filter(s => new Date(s.start_time).getTime() > now)
      .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime());

    this.nextShift = futureShifts.length > 0 ? futureShifts[0] : null;
  }

  private getAssignedShifts(shifts: WorkerShift[]) {
    if (!shifts) return [];
    return shifts.filter(shift => shift.status?.toLowerCase() === 'assigned');
  }

  calculateWeeklyAssigned(shifts: WorkerShift[]) {
    const startOfWeek = this.getWeekStartFromOffset();
    const endOfWeek = new Date(startOfWeek);
    endOfWeek.setDate(startOfWeek.getDate() + 6);
    endOfWeek.setHours(23, 59, 59, 999);

    this.weekRangeLabel = `${startOfWeek.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${endOfWeek.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}`;

    const dayLabels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
    this.weeklyDays = dayLabels.map((label, index) => {
      const dayDate = new Date(startOfWeek);
      dayDate.setDate(startOfWeek.getDate() + index);
      return { label, date: dayDate, hours: 0, percent: 0 };
    });

    this.weeklyTotal = 0;

    shifts.forEach(shift => {
      if (!shift.start_time || !shift.end_time) return;

      const status = shift.status?.toLowerCase();
      if (status && status !== 'assigned' && status !== 'completed') return;

      const start = new Date(shift.start_time);
      const end = new Date(shift.end_time);

      if (start < startOfWeek || start > endOfWeek) return;

      const durationHours = (end.getTime() - start.getTime()) / (1000 * 60 * 60);
      if (durationHours <= 0) return;

      const dayIndex = (start.getDay() + 6) % 7;
      this.weeklyDays[dayIndex].hours += durationHours;
      this.weeklyTotal += durationHours;
    });

    this.weeklyDays = this.weeklyDays.map(day => ({
      ...day,
      percent: Math.min(100, (day.hours / this.dailyTarget) * 100)
    }));

    this.weeklyProgress = Math.min(100, (this.weeklyTotal / this.weeklyTarget) * 100);
  }

  private getStartOfWeek(date: Date) {
    const start = new Date(date);
    const day = start.getDay();
    const diff = (day + 6) % 7;
    start.setDate(start.getDate() - diff);
    start.setHours(0, 0, 0, 0);
    return start;
  }

  private getWeekStartFromOffset() {
    const base = this.getStartOfWeek(new Date());
    base.setDate(base.getDate() + this.weekOffset * 7);
    return base;
  }

  previousWeek() {
    this.weekOffset -= 1;
    this.refreshWeeklyView();
  }

  nextWeek() {
    this.weekOffset += 1;
    this.refreshWeeklyView();
  }

  private refreshWeeklyView() {
    const safeData = this.latestShifts || [];
    this.calculateWeeklyAssigned(safeData);
  }

  private updateDateTime() {
    const now = new Date();
    this.currentDateTime = new Intl.DateTimeFormat('en-US', {
      weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
    }).format(now);
  }

  toggleProfileMenu() {
    this.isProfileMenuOpen = !this.isProfileMenuOpen;
  }

  closeProfileMenu() {
    this.isProfileMenuOpen = false;
  }

  logout() {
    this.authService.logout();
    this.closeProfileMenu();
    this.router.navigate(['/login']);
  }
}
