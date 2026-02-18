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
export class WorkerDashboardComponent implements OnInit{

  // --- HARDCODED WORKER ID (Remove when Login is built) ---
  currentWorkerId = 3;
  currentWorkerName = 'Worker3';

  currentDateTime = '';

  myShifts$!: Observable<WorkerShift[]>;
  nextShift: WorkerShift | null = null;

  totalHours = 0;
  estEarnings = 0;
  hourlyWage = 15;

  constructor(private workerService: WorkerService) { }

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

    // 2. Assign Observable with a 'tap' side-effect
    // The 'tap' operator lets us run code (calculate stats) 
    // without stopping the data from going to the HTML.
    this.myShifts$ = this.workerService.getMyShifts(this.currentWorkerId).pipe(
      tap((data) => {
        console.log("Stream received data:", data);
        this.calculateStats(data);
        this.findNextShift(data);
      })
    );
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

    // Filter for future shifts and sort by start time
    const futureShifts = shifts
      .filter(s => new Date(s.start_time).getTime() > now)
      .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime());

    this.nextShift = futureShifts.length > 0 ? futureShifts[0] : null;
  }

  private updateDateTime() {
    const now = new Date();
    this.currentDateTime = new Intl.DateTimeFormat('en-US', {
      weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
    }).format(now);
  }
}