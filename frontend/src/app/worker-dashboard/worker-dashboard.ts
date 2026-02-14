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
      
      <mat-toolbar>
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
  /* 1. Reset the container to let the global gradient show through */
  .dashboard-container { 
    min-height: 100vh; 
    /* No background-color here; we want the purple/blue global gradient */
  }

  /* 2. Style the Toolbar to look premium */
  mat-toolbar {
    background: rgba(255, 255, 255, 0.9); /* Semi-transparent white */
    color: #333; /* Dark text for contrast */
    box-shadow: 0 2px 10px rgba(0,0,0,0.1); /* Subtle shadow */
    backdrop-filter: blur(5px);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .spacer { flex: 1 1 auto; }

  /* 3. The Content Area */
  .content { 
    padding: 40px; 
    max-width: 1200px; 
    margin: 0 auto; 
  }

  h1 { 
    color: white; /* White text looks great on the purple gradient */
    font-size: 2rem;
    font-weight: 300;
    margin-bottom: 20px;
    text-shadow: 0 2px 4px rgba(0,0,0,0.2); /* Make text readable */
  }

  /* 4. The Grid Card - The "Glass Sheet" */
  .grid-card { 
    padding: 0; /* Remove padding so grid fills the card */
    overflow: hidden; /* Round the corners of the grid */
    
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.2);
  }
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