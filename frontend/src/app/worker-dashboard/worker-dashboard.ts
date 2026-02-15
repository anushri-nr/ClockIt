import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AgGridAngular } from 'ag-grid-angular'; 
import { ColDef, CellClassParams } from 'ag-grid-community'; 
import { ModuleRegistry, AllCommunityModule } from 'ag-grid-community';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';

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
  template: `
    <div class="dashboard-wrapper">
      
      <mat-toolbar class="custom-toolbar">
        <span class="brand">ClockIt</span>
        <span class="spacer"></span>
        <div class="user-profile">
          <span>Hello, Varshith</span>
          <button mat-icon-button>
            <mat-icon>account_circle</mat-icon>
          </button>
        </div>
      </mat-toolbar>

      <div class="content-container">
        
        <div class="stats-grid">
          <mat-card class="stat-card highlight-card">
            <mat-card-header>
              <mat-card-subtitle>Next Shift</mat-card-subtitle>
              <mat-card-title>Tomorrow, 9:00 AM</mat-card-title>
            </mat-card-header>
            <mat-card-content>
              <p>Main Store • Cashier</p>
            </mat-card-content>
            <mat-card-actions>
              <button mat-button class="btn-inverse">View Details</button>
            </mat-card-actions>
          </mat-card>

          <mat-card class="stat-card">
            <div class="icon-header">
              <mat-icon class="green-icon">schedule</mat-icon>
              <span>Weekly Hours</span>
            </div>
            <div class="big-number">32.5</div>
            <div class="sub-text">Target: 40 hrs</div>
          </mat-card>

          <mat-card class="stat-card">
             <div class="icon-header">
              <mat-icon class="orange-icon">payments</mat-icon>
              <span>Est. Earnings</span>
            </div>
            <div class="big-number">$487.50</div>
            <div class="sub-text">Pending Approval</div>
          </mat-card>
        </div>

        <div class="grid-section">
          <h2>My Schedule</h2>
          <mat-card class="grid-container">
            <ag-grid-angular
              class="ag-theme-quartz"
              style="height: 400px; width: 100%;"
              [rowData]="shifts"
              [columnDefs]="colDefs">
            </ag-grid-angular>
          </mat-card>
        </div>

      </div>
    </div>
  `,
  styles: [`
    /* LAYOUT & STRUCTURE */
    .dashboard-wrapper {
      min-height: 100vh;
      background: radial-gradient(circle at 10% 10%, #fef3e6 0%, transparent 40%),
        radial-gradient(circle at 90% 0%, #e8f6f2 0%, transparent 45%),
        linear-gradient(180deg, #f7f2ea 0%, #f1ebe1 100%);
      position: relative;
      overflow: hidden;
    }

    .dashboard-wrapper::before {
      content: '';
      position: absolute;
      inset: -20% -10% auto auto;
      width: 420px;
      height: 420px;
      background: conic-gradient(from 120deg, rgba(244, 162, 97, 0.2), rgba(15, 95, 92, 0.2), transparent 65%);
      border-radius: 50%;
      filter: blur(10px);
      z-index: 0;
    }

    .dashboard-wrapper > * {
      position: relative;
      z-index: 1;
    }

    .content-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 30px 20px;
    }

    /* TOOLBAR STYLING */
    .custom-toolbar {
      background: linear-gradient(135deg, #0f5f5c 0%, #144e4b 45%, #1f2d2d 100%);
      color: white;
      box-shadow: 0 6px 18px rgba(15, 95, 92, 0.25);
      border-radius: 16px;
      margin: 16px auto 24px;
      padding: 6px 12px;
      width: min(1200px, calc(100% - 40px));
      box-sizing: border-box;
    }
    
    .brand {
      font-weight: 700;
      letter-spacing: 1px;
      font-size: 1.2rem;
    }

    .user-profile {
      display: flex;
      align-items: center;
      gap: 10px;
      font-size: 0.9rem;
    }

    .spacer { flex: 1 1 auto; }

    /* WIDGET GRID (The "New Functionality") */
    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 20px;
      margin-bottom: 30px;
    }

    .stat-card {
      padding: 15px;
      background: #fffdf9;
      border: 1px solid rgba(31, 35, 40, 0.08);
      box-shadow: 0 12px 30px rgba(31, 35, 40, 0.08);
    }

    /* Special styling for the "Next Shift" card */
    .highlight-card {
      background: linear-gradient(135deg, #0f5f5c 0%, #144e4b 100%);
      color: white;
    }
    
    .highlight-card mat-card-subtitle { color: rgba(255,255,255, 0.7); }
    .highlight-card mat-card-title { color: white; font-size: 1.5rem; margin-bottom: 10px; }
    
    .btn-inverse {
      background: rgba(255,255,255,0.2);
      color: white;
    }

    /* Standard Stat Cards */
    .icon-header {
      display: flex;
      align-items: center;
      gap: 10px;
      color: #666;
      font-size: 0.9rem;
      text-transform: uppercase;
      letter-spacing: 1px;
    }

    .big-number {
      font-size: 2.5rem;
      font-weight: 700;
      margin: 15px 0;
      color: #1f2328;
    }

    .sub-text { color: #888; font-size: 0.9rem; }
    
    .green-icon { color: #0f5f5c; }
    .orange-icon { color: #f4a261; }

    /* GRID SECTION */
    h2 {
      font-weight: 600;
      color: #1f2328;
      margin-bottom: 15px;
    }

    .grid-container {
      padding: 0; 
      overflow: hidden;
      background: #fffdf9;
      border: 1px solid rgba(31, 35, 40, 0.08);
      box-shadow: 0 12px 30px rgba(31, 35, 40, 0.08);
    }

    .ag-theme-quartz {
      --ag-background-color: #fffdf9;
      --ag-header-background-color: #f4efe7;
      --ag-odd-row-background-color: #fffaf3;
      --ag-row-hover-color: rgba(15, 95, 92, 0.06);
      --ag-border-color: rgba(31, 35, 40, 0.12);
      --ag-header-foreground-color: #1f2328;
      --ag-foreground-color: #1f2328;
    }
  `]
})
export class WorkerDashboardComponent {
  
  shifts: Shift[] = [
    { date: '2026-02-08', startTime: '09:00', endTime: '17:00', status: 'Assigned', hourlyWage: 15 },
    { date: '2026-02-09', startTime: '09:00', endTime: '17:00', status: 'Assigned', hourlyWage: 15 },
    { date: '2026-02-10', startTime: '12:00', endTime: '20:00', status: 'Assigned', hourlyWage: 15.5 },
    { date: '2026-02-12', startTime: '09:00', endTime: '17:00', status: 'Completed', hourlyWage: 15 },
  ];

  colDefs: ColDef[] = [
     { field: 'date', flex: 1 },
     { field: 'startTime', headerName: 'Start', flex: 1 },
     { field: 'endTime', headerName: 'End', flex: 1 },
     { field: 'status', flex: 1, cellStyle: (params: CellClassParams) => {
        return { 
          color: params.value === 'Assigned' ? '#0f5f5c' : '#666', 
          fontWeight: 'bold' 
        };
     }},
     { field: 'hourlyWage', headerName: 'Wage ($/hr)', flex: 1 },
  ];
}
