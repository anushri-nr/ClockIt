import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AgGridAngular } from 'ag-grid-angular'; // The Grid Component
import { ColDef } from 'ag-grid-community'; // The Type Definition
import { ModuleRegistry, AllCommunityModule } from 'ag-grid-community';

// Register all Community features
ModuleRegistry.registerModules([AllCommunityModule]);

// Our Contract (The Interface we agreed on)
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
  imports: [CommonModule, AgGridAngular],
  template: `
    <div style="padding: 20px;">
      <h1>My Schedule</h1>
      <ag-grid-angular
        class="ag-theme-quartz"
        style="height: 400px; width: 100%;"
        [rowData]="shifts"
        [columnDefs]="colDefs">
      </ag-grid-angular>
    </div>
  `
})
export class WorkerDashboardComponent {
  // 1. The Mock Data (The "Meat")
  shifts: Shift[] = [
    { date: '2026-02-08', startTime: '09:00', endTime: '17:00', status: 'Assigned', hourlyWage: 15 },
    { date: '2026-02-09', startTime: '09:00', endTime: '17:00', status: 'Assigned', hourlyWage: 15 },
  ];

  // 2. The Column Definitions (The "Recipe")
  // YOUR JOB: Define the columns to display the data above.
  // Syntax: { field: 'fieldName' }
  colDefs: ColDef[] = [
     // WRITE YOUR CODE HERE
     // Example: { field: 'date' },
     { field: 'date' },
     { field: 'startTime' },
     { field: 'endTime' },
     { field: 'status' },
     { field: 'hourlyWage' },

  ];
}