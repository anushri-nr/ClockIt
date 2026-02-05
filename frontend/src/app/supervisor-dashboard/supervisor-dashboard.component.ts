import { Component } from '@angular/core';
import { NgIf } from '@angular/common';

@Component({
  selector: 'app-supervisor-dashboard',
  standalone: true,
  imports: [NgIf],
  templateUrl: './supervisor-dashboard.component.html',
  styleUrl: './supervisor-dashboard.component.scss'
})
export class SupervisorDashboardComponent {
  isShiftModalOpen = false;
  isAnnouncementModalOpen = false;
  isOpenShiftsModalOpen = false;
  isAssignWorkerModalOpen = false;
  isStatusSnapshotOpen = false;
  scheduleView: 'today' | 'week' = 'today';

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
