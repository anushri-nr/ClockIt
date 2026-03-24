import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ShiftService } from '../services/shift.service';
import { SupervisorService } from '../services/supervisor.service';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-supervisor-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule, MatToolbarModule, MatIconModule, MatButtonModule],
  templateUrl: './supervisor-dashboard.component.html',
  styleUrl: './supervisor-dashboard.component.scss'
})
export class SupervisorDashboardComponent implements OnInit {

  // Dynamic User Data (Replaces Hardcoded Values)
  currentSupervisorId: number = 0;
  currentCompanyId: number = 0;
  currentSupervisorName: string = '';

  // View State
  currentDateTime = '';
  scheduleView: 'today' | 'week' = 'today';

  // Modal Flags
  isShiftModalOpen = false;
  isAnnouncementModalOpen = false;
  isOpenShiftsModalOpen = false;
  isAssignWorkerModalOpen = false;
  isStatusSnapshotOpen = false;

  // Create Shift Data
  createShift = { startTime: '', endTime: '', createdBy: this.currentSupervisorId.toString() };
  createdAt = '';
  isSubmittingShift = false;

  // NEW: Error Messages
  createShiftError = '';
  assignShiftError = '';

  // Assign Worker Data
  selectedShiftId = '';
  selectedWorkerId = '';
  availableWorkers: any[] = [];
  isLoadingWorkers = false;

  // Mock Data (Placeholder until 'Get Shifts' API is ready)
  shifts: any[] = [];
  openShifts: any[] = [];
  pendingApprovalShifts: any[] = [];
  isLoadingOpenShifts = false;
  isLoadingPendingShifts = false;

  constructor(
    private shiftService: ShiftService,
    private supervisorService: SupervisorService,
    private authService: AuthService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit() {
    this.updateDateTime();
    this.createdAt = new Date().toISOString();
    // 1. Ask the single source of truth for the user
    const user = this.authService.currentUserValue;

    // 2. Safely assign the variables
    if (user && user.id) {
      this.currentSupervisorId = user.id;
      this.currentSupervisorName = user.name || 'Supervisor';
      this.currentCompanyId = user.company_id || 1;

      console.log(`Dashboard initialized for User ID: ${this.currentSupervisorId}`);
      
      // 3. Load the shifts using a guaranteed valid ID!
      this.loadShifts(); 
    } else {
      console.error('CRITICAL: No valid user state found. Cannot load shifts.');
    }
  }

  loadShifts() {
    this.supervisorService.getShifts(this.currentSupervisorId).subscribe({
      next: (data) => {
        console.log('Shifts Loaded:', data);
        this.shifts = data;
      },
      error: (err) => console.error('Failed to load shifts', err)
    });
  }

  // --- LOGIC ---

  get totalHours(): string {
    if (!this.createShift.startTime || !this.createShift.endTime) return '0.0';
    const start = new Date(this.createShift.startTime);
    const end = new Date(this.createShift.endTime);
    let diffMs = end.getTime() - start.getTime();
    if (diffMs <= 0) diffMs += 24 * 60 * 60 * 1000; // Handle midnight crossing
    return (diffMs / 36e5).toFixed(1);
  }

  submitCreateShift(): void {
    if (this.isSubmittingShift) return;
    
    // Reset Error
    this.createShiftError = '';

    this.createShift.createdBy = this.currentSupervisorId.toString();

    if (!this.createShift.startTime || !this.createShift.endTime) {
      alert('Please fill start time and end times.');
      return;
    }

    // Validate End Time > Start Time
    const start = new Date(this.createShift.startTime).getTime();
    const end = new Date(this.createShift.endTime).getTime();

    if (end <= start) {
      this.createShiftError = 'End time must be after start time.';
      return;
    }

    this.isSubmittingShift = true;

    // Use the Service, not direct HTTP
    this.shiftService.createShift(this.createShift).subscribe({
      next: () => {
        console.log("Shift created successfully");

        this.isSubmittingShift = false;
        this.closeShiftModal();

        this.cdr.detectChanges();
        
        // Reset form
        this.createShift = { startTime: '', endTime: '', createdBy: this.currentSupervisorId.toString() };
        this.loadShifts();
      },
      error: (err) => {
        console.error(err);
        this.isSubmittingShift = false;
        this.createShiftError = 'Failed to create shift. Please try again.';
        this.cdr.detectChanges();
      }
    });
  }

  assignWorkerToShift(): void {
    this.assignShiftError = '';

    if (!this.selectedShiftId || !this.selectedWorkerId) {
      alert('Please select a shift and a worker.');
      return;
    }

    // Convert strings to numbers for the API
    const shiftIdNum = Number(this.selectedShiftId);
    const workerIdNum = Number(this.selectedWorkerId);

    this.supervisorService.assignWorker(shiftIdNum, workerIdNum, this.currentSupervisorId).subscribe({
      next: () => {
        console.log("Worker assigned successfully");

        this.closeAssignWorkerModal();

        this.cdr.detectChanges();
        
        // Refresh the grid to show the new assignment
        this.loadShifts(); 
      },
      error: (err) => {
        console.error('Assignment failed', err);
        this.assignShiftError = 'Failed to assign worker. Please try again.';
        this.cdr.detectChanges();
      }
    });
  }

  onShiftSelected() {
    // 1. Reset current worker list
    this.availableWorkers = [];
    this.selectedWorkerId = '';
    this.assignShiftError = ''; // Clear errors on selection change

    if(!this.selectedShiftId) return;

    // 2. Find the full shift object to get its date
    // We cast to 'any' to avoid type errors with snake_case
    const shift: any = this.shifts.find((s: any) => s.id == this.selectedShiftId);

    if(shift) {
      this.isLoadingWorkers = true;

      // 3. Extract the date (YYYY-MM-DD) from the shift's start time
      // Example: "2026-02-18T09:00:00Z" -> "2026-02-18"
      const dateStr = shift.start_time.split('T')[0];

      // 4. service to get workers for THAT date
      this.supervisorService.getAvailableWorkers(dateStr, this.currentCompanyId).subscribe({
        next: (data) => {
          this.availableWorkers = data;
          this.isLoadingWorkers = false;
          console.log(`Loaded ${data.length} workers for date: ${dateStr}`);
        },
        error: (err) => {
          console.error('Failed to load workers for selected shift', err);
          this.isLoadingWorkers = false;
        }
      });
    }

  }

  // --- MODAL CONTROLS ---

  private resetModals() {
    this.isShiftModalOpen = false;
    this.isAnnouncementModalOpen = false;
    this.isOpenShiftsModalOpen = false;
    this.isAssignWorkerModalOpen = false;
    this.isStatusSnapshotOpen = false;
    this.createShiftError = '';
    this.assignShiftError = '';
  }

  openShiftModal() {
    this.resetModals();
    this.isShiftModalOpen = true;
  }

  closeShiftModal() {
    this.isShiftModalOpen = false;
  }

  openAssignWorkerModal() {
  this.resetModals();
  this.isAssignWorkerModalOpen = true;
  
  this.selectedShiftId = '';
  this.selectedWorkerId = '';
  this.availableWorkers = [];
}

  closeAssignWorkerModal() {
    this.isAssignWorkerModalOpen = false;
  }

  // Simple toggles
  openAnnouncementModal() { this.resetModals(); this.isAnnouncementModalOpen = true; }
  closeAnnouncementModal() { this.isAnnouncementModalOpen = false; }

  // Fetch 'Unassigned' shifts when opening the Open Shifts Modal
  openOpenShiftsModal() { 
    this.resetModals(); 
    this.isOpenShiftsModalOpen = true; 
    this.isLoadingOpenShifts = true;

    this.supervisorService.getShifts(this.currentSupervisorId, 'Unassigned').subscribe({
      next: (data) => {
        this.openShifts = data || []; 
        this.isLoadingOpenShifts = false;
        this.cdr.detectChanges(); // Force the UI to refresh instantly
      },
      error: (err) => {
        console.error('Failed to load open shifts', err);
        this.openShifts = []; // Clear array
        this.isLoadingOpenShifts = false; // Stop spinner
        this.cdr.detectChanges(); // Force UI to update
      }
    });
  }
  closeOpenShiftsModal() { this.isOpenShiftsModalOpen = false; }

  // Fetch 'Requested' shifts when opening the Status Snapshot (Pending Approval) Modal
  openStatusSnapshot() { 
    this.resetModals(); 
    this.isStatusSnapshotOpen = true; 
    this.isLoadingPendingShifts = true;

    this.supervisorService.getShifts(this.currentSupervisorId, 'Requested').subscribe({
      next: (data) => {
        this.pendingApprovalShifts = data || []; 
        this.isLoadingPendingShifts = false;
        this.cdr.detectChanges(); // Force the UI to refresh instantly
      },
      error: (err) => {
        console.error('Failed to load pending approval shifts', err);
        this.pendingApprovalShifts = []; // Clear array
        this.isLoadingPendingShifts = false; // Stop spinner
        this.cdr.detectChanges(); // Force UI to update
      }
    });
  }
  closeStatusSnapshot() { this.isStatusSnapshotOpen = false; }

  setScheduleView(view: 'today' | 'week') {
    this.scheduleView = view;
  }

  private updateDateTime() {
    const now = new Date();
    this.currentDateTime = new Intl.DateTimeFormat('en-US', {
      weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
    }).format(now);
  }
}