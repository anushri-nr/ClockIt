import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AgGridAngular } from 'ag-grid-angular'; 
import { ColDef } from 'ag-grid-community'; 
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
 
  .dashboard-wrapper {
    min-height: 100vh;
    padding-bottom: 40px; 
    
    
    background: 
      radial-gradient(circle at 10% 10%, rgba(254, 243, 230, 1) 0%, transparent 40%),
      radial-gradient(circle at 90% 0%, rgba(232, 246, 242, 1) 0%, transparent 45%),
      linear-gradient(180deg, #f7f2ea 0%, #f1ebe1 100%);
      
    position: relative;
    overflow-x: hidden;
  }

 
  .dashboard-wrapper::before {
    content: '';
    position: absolute;
    top: -100px;
    right: -100px;
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(15, 95, 92, 0.08) 0%, transparent 70%);
    border-radius: 50%;
    z-index: 0;
    pointer-events: none;
  }


  mat-toolbar {
    background: linear-gradient(135deg, #0f5f5c 0%, #144e4b 45%, #1f2d2d 100%);
    color: white;
    width: calc(100% - 40px);
    max-width: 1200px;
    margin: 20px auto 30px; 
    border-radius: 16px;
    box-shadow: 0 6px 18px rgba(15, 95, 92, 0.25);
    position: sticky;
    top: 20px; 
    z-index: 100;
  }
  
  .content-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
    position: relative;
    z-index: 1;
  }

 
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 20px;
    margin-bottom: 30px;
  }

  .stat-card {
    padding: 20px;
    border-radius: 20px !important; 
    
    
    background-color: #fffdf9 !important; 
    
    
    border: 1px solid rgba(15, 95, 92, 0.08);
    box-shadow: 0 14px 30px rgba(31, 35, 40, 0.08) !important;
  }

 
  .highlight-card {
    background: linear-gradient(135deg, #0f5f5c 0%, #144e4b 100%) !important;
    color: white;
    box-shadow: 0 14px 30px rgba(15, 95, 92, 0.25) !important;
  }
  
  .highlight-card mat-card-title { color: white; font-size: 1.8rem; margin-bottom: 5px; }
  .highlight-card mat-card-subtitle { color: rgba(255,255,255, 0.8); }

  
  h2 {
    font-weight: 700;
    color: #1f2328;
    margin-bottom: 15px;
    font-size: 1.2rem;
    letter-spacing: 0.05em;
  }

  .grid-container {
    padding: 25px;
    border-radius: 24px !important;
    

    background-color: #fffdf9 !important;
    
    border: 1px solid rgba(15, 95, 92, 0.08);
    box-shadow: 0 14px 30px rgba(31, 35, 40, 0.08) !important;
  }


  .brand { font-weight: 700; letter-spacing: 1px; font-size: 1.1rem; }
  .spacer { flex: 1 1 auto; }
  
  .big-number {
    font-size: 2.8rem;
    font-weight: 700;
    margin: 10px 0;
    color: #1f2328;
  }
  
  .sub-text { color: #5a6472; font-size: 0.9rem; font-weight: 500; }
  
  .icon-header {
    display: flex; align-items: center; gap: 8px;
    color: #5a6472; font-size: 0.85rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.08em;
  }

  button mat-icon { color: white; }
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
     { field: 'status', flex: 1, cellStyle: (params) => {
        return { 
          color: params.value === 'Assigned' ? '#0f5f5c' : '#666', 
          fontWeight: 'bold' 
        };
     }},
     { field: 'hourlyWage', headerName: 'Wage ($/hr)', flex: 1 },
  ];
}