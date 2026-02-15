import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ShiftService } from '../services/shift';
import { SupervisorService } from '../services/supervisor';

@Component({
  selector: 'app-supervisor-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule, MatToolbarModule, MatIconModule, MatButtonModule],
  templateUrl: './supervisor-dashboard.html',
  styleUrl: './supervisor-dashboard.scss'
})
export class SupervisorDashboardComponent implements OnInit {
  
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
  createShift = { startTime: '', endTime: '', createdBy: '' };
  createdAt = '';
  isSubmittingShift = false;

  // Assign Worker Data
  selectedShiftId = '';
  selectedWorkerId = '';
  availableWorkers: any[] = [];
  isLoadingWorkers = false;

  // Mock Data (Placeholder until 'Get Shifts' API is ready)
  mockShifts = [
    { id: 'shift-001', startTime: '2026-02-15T09:00', endTime: '2026-02-15T13:00', createdBy: 'Supervisor A', createdAt: '2026-02-15T08:30:00Z' },
    { id: 'shift-002', startTime: '2026-02-15T14:00', endTime: '2026-02-15T18:00', createdBy: 'Supervisor B', createdAt: '2026-02-15T09:10:00Z' }
  ];

  constructor(
    private shiftService: ShiftService,
    private supervisorService: SupervisorService
  ) {}

  ngOnInit() {
    this.updateDateTime();
    this.createdAt = new Date().toISOString();
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
    
    if (!this.createShift.startTime || !this.createShift.endTime || !this.createShift.createdBy) {
      alert('Please fill all fields before submitting.');
      return;
    }

    this.isSubmittingShift = true;

    // Use the Service, not direct HTTP
    this.shiftService.createShift(this.createShift).subscribe({
      next: () => {
        alert('Shift created successfully!');
        this.isSubmittingShift = false;
        this.closeShiftModal();
        // Reset form
        this.createShift = { startTime: '', endTime: '', createdBy: '' };
      },
      error: (err) => {
        console.error(err);
        this.isSubmittingShift = false;
        alert('Failed to create shift.');
      }
    });
  }

  assignWorkerToShift(): void {
    if (!this.selectedShiftId || !this.selectedWorkerId) {
      alert('Please select a shift and a worker.');
      return;
    }
    // TODO: Connect to POST /api/shifts/assign when ready
    console.log(`Assigning worker ${this.selectedWorkerId} to shift ${this.selectedShiftId}`);
    this.closeAssignWorkerModal();
  }

  private updateDateTime() {
    const now = new Date();
    this.currentDateTime = new Intl.DateTimeFormat('en-US', {
      weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
    }).format(now);
  }

  // --- MODAL CONTROLS ---

  private resetModals() {
    this.isShiftModalOpen = false;
    this.isAnnouncementModalOpen = false;
    this.isOpenShiftsModalOpen = false;
    this.isAssignWorkerModalOpen = false;
    this.isStatusSnapshotOpen = false;
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
    
    // FETCH WORKERS ON OPEN
    this.isLoadingWorkers = true;
    this.supervisorService.getAvailableWorkers().subscribe({
      next: (data) => {
        this.availableWorkers = data;
        this.isLoadingWorkers = false;
      },
      error: (err) => {
        console.error(err);
        this.isLoadingWorkers = false;
      }
    });
  }

  closeAssignWorkerModal() {
    this.isAssignWorkerModalOpen = false;
  }

  // Simple toggles
  openAnnouncementModal() { this.resetModals(); this.isAnnouncementModalOpen = true; }
  closeAnnouncementModal() { this.isAnnouncementModalOpen = false; }
  
  openOpenShiftsModal() { this.resetModals(); this.isOpenShiftsModalOpen = true; }
  closeOpenShiftsModal() { this.isOpenShiftsModalOpen = false; }
  
  openStatusSnapshot() { this.resetModals(); this.isStatusSnapshotOpen = true; }
  closeStatusSnapshot() { this.isStatusSnapshotOpen = false; }

  setScheduleView(view: 'today' | 'week') {
    this.scheduleView = view;
  }
}