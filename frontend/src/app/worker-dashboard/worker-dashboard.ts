import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AgGridAngular } from 'ag-grid-angular'; 
import { ColDef } from 'ag-grid-community'; 
import { ModuleRegistry, AllCommunityModule } from 'ag-grid-community';

// MATERIAL IMPORTS
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';

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
    MatCardModule
  ],
  template: `
    <div class="dashboard-container">
      
      <mat-toolbar color="primary">
        <span>ClockIt: Worker Portal</span>
        <span class="spacer"></span>
        <button mat-icon-button>
          <mat-icon>account_circle</mat-icon>
        </button>
      </mat-toolbar>

      <div class="content">
        <h1>My Schedule</h1>
        
        <mat-card class="grid-card">
          <ag-grid-angular
            class="ag-theme-quartz"
            style="height: 500px; width: 100%;"
            [rowData]="shifts"
            [columnDefs]="colDefs">
          </ag-grid-angular>
        </mat-card>
      </div>

    </div>
  `,
  styles: [`
    .dashboard-container { background-color: #f5f5f5; min-height: 100vh; }
    .content { padding: 20px; max-width: 1200px; margin: 0 auto; }
    .spacer { flex: 1 1 auto; }
    .grid-card { padding: 20px; }
    h1 { margin-bottom: 20px; color: #333; }
  `]
})
export class WorkerDashboardComponent {
  // Mock Data
  shifts: Shift[] = [
    { date: '2026-02-08', startTime: '09:00', endTime: '17:00', status: 'Assigned', hourlyWage: 15 },
    { date: '2026-02-09', startTime: '09:00', endTime: '17:00', status: 'Assigned', hourlyWage: 15 },
  ];

  colDefs: ColDef[] = [
     { field: 'date', flex: 1 },
     { field: 'startTime', headerName: 'Start', flex: 1 },
     { field: 'endTime', headerName: 'End', flex: 1 },
     { field: 'status', flex: 1, cellStyle: { 'font-weight': 'bold', 'color': 'green' } },
     { field: 'hourlyWage', headerName: 'Wage ($/hr)', flex: 1 },
  ];
}