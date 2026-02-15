import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-supervisor-dashboard',
  standalone: true,
  imports: [CommonModule, FormsModule, MatToolbarModule, MatIconModule, MatButtonModule],
  templateUrl: './supervisor-dashboard.component.html',
  styleUrl: './supervisor-dashboard.component.scss'
})
export class SupervisorDashboardComponent implements OnInit {
  currentDateTime = '';
  createShift = {
    startTime: '',
    endTime: '',
    createdBy: ''
  };
  createdAt = '';
  isSubmittingShift = false;
  selectedShiftId = '';
  selectedWorkerId = '';

  mockShifts = [
    {
      id: 'shift-001',
      startTime: '2026-02-15T09:00',
      endTime: '2026-02-15T13:00',
      createdBy: 'Supervisor A',
      createdAt: '2026-02-15T08:30:00Z'
    },
    {
      id: 'shift-002',
      startTime: '2026-02-15T14:00',
      endTime: '2026-02-15T18:00',
      createdBy: 'Supervisor B',
      createdAt: '2026-02-15T09:10:00Z'
    },
    {
      id: 'shift-003',
      startTime: '2026-02-16T08:00',
      endTime: '2026-02-16T12:00',
      createdBy: 'Supervisor C',
      createdAt: '2026-02-15T10:05:00Z'
    }
  ];

  mockWorkers = [
    { id: 'W-1001', name: 'Maria L.' },
    { id: 'W-1002', name: 'Devon P.' },
    { id: 'W-1003', name: 'Jalen K.' },
    { id: 'W-1004', name: 'Nora K.' }
  ];
  isShiftModalOpen = false;
  isAnnouncementModalOpen = false;
  isOpenShiftsModalOpen = false;
  isAssignWorkerModalOpen = false;
  isStatusSnapshotOpen = false;
  scheduleView: 'today' | 'week' = 'today';

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.updateDateTime();
    this.createdAt = new Date().toISOString();
  }

  private updateDateTime() {
    const now = new Date();
    this.currentDateTime = new Intl.DateTimeFormat('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    }).format(now);
  }

  get totalHours(): string {
    if (!this.createShift.startTime || !this.createShift.endTime) {
      return '0.0';
    }
    const start = new Date(this.createShift.startTime);
    const end = new Date(this.createShift.endTime);
    let diffMs = end.getTime() - start.getTime();
    if (Number.isNaN(diffMs) || diffMs <= 0) {
      // If end is earlier than start, assume it crosses midnight into the next day.
      if (!Number.isNaN(diffMs)) {
        diffMs += 24 * 60 * 60 * 1000;
      }
    }
    if (Number.isNaN(diffMs) || diffMs <= 0) {
      return '0.0';
    }
    return (diffMs / 36e5).toFixed(1);
  }

  submitCreateShift(): void {
    if (this.isSubmittingShift) return;
    if (!this.createShift.startTime || !this.createShift.endTime || !this.createShift.createdBy) {
      alert('Please fill all fields before submitting.');
      return;
    }

    this.isSubmittingShift = true;
    const payload = {
      startTime: this.createShift.startTime,
      endTime: this.createShift.endTime,
      createdBy: this.createShift.createdBy,
      createdAt: this.createdAt,
      totalHours: Number(this.totalHours)
    };

    this.http.post('/api/shifts', payload).subscribe({
      next: () => {
        this.isSubmittingShift = false;
        this.closeShiftModal();
      },
      error: () => {
        this.isSubmittingShift = false;
        alert('Failed to create shift. Please try again.');
      }
    });
  }

  assignWorkerToShift(): void {
    if (!this.selectedShiftId || !this.selectedWorkerId) {
      alert('Please select a shift and a worker.');
      return;
    }
    // Mocked until API is ready.
    this.closeAssignWorkerModal();
  }

  openShiftModal() {
    this.isShiftModalOpen = true;
    this.isAnnouncementModalOpen = false;
    this.isOpenShiftsModalOpen = false;
    this.isAssignWorkerModalOpen = false;
    this.isStatusSnapshotOpen = false;
  }

  closeShiftModal() {
    this.isShiftModalOpen = false;
  }

  openAnnouncementModal() {
    this.isShiftModalOpen = false;
    this.isAnnouncementModalOpen = true;
    this.isOpenShiftsModalOpen = false;
    this.isAssignWorkerModalOpen = false;
    this.isStatusSnapshotOpen = false;
  }

  openOpenShiftsModal() {
    this.isShiftModalOpen = false;
    this.isAnnouncementModalOpen = false;
    this.isOpenShiftsModalOpen = true;
    this.isAssignWorkerModalOpen = false;
    this.isStatusSnapshotOpen = false;
  }

  closeOpenShiftsModal() {
    this.isOpenShiftsModalOpen = false;
  }

  openAssignWorkerModal() {
    this.isShiftModalOpen = false;
    this.isAnnouncementModalOpen = false;
    this.isOpenShiftsModalOpen = false;
    this.isAssignWorkerModalOpen = true;
    this.isStatusSnapshotOpen = false;
  }

  closeAssignWorkerModal() {
    this.isAssignWorkerModalOpen = false;
  }

  openStatusSnapshot() {
    this.isShiftModalOpen = false;
    this.isAnnouncementModalOpen = false;
    this.isOpenShiftsModalOpen = false;
    this.isAssignWorkerModalOpen = false;
    this.isStatusSnapshotOpen = true;
  }

  closeStatusSnapshot() {
    this.isStatusSnapshotOpen = false;
  }

  closeAnnouncementModal() {
    this.isAnnouncementModalOpen = false;
  }

  setScheduleView(view: 'today' | 'week') {
    this.scheduleView = view;
  }
}
