import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-shift-list',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  template: `
    <div *ngIf="isLoading" style="padding: 10px; color: #5a6472;">
      {{ loadingMessage }}
    </div>
    
    <div *ngIf="!isLoading && shifts.length === 0" style="padding: 10px; color: #5a6472;">
      {{ emptyMessage }}
    </div>

    <ul class="list compact" *ngIf="!isLoading && shifts.length > 0">
      <li *ngFor="let shift of shifts" style="margin-bottom: 12px; padding: 16px; border: 1px solid #e2e8f0; border-radius: 8px;">
        <div style="display: flex; justify-content: space-between; align-items: center; width: 100%;">
          
          <div>
            <div class="request-title" style="font-weight: 600;">
              <ng-container *ngIf="showWorkerName; else timeOnly">
                Worker: {{ shift.assigned_to || 'Unknown' }} requested a shift
              </ng-container>
              <ng-template #timeOnly>
                {{ shift.start_time | date:'MMM d, y, h:mm a' }} → {{ shift.end_time | date:'h:mm a' }}
              </ng-template>
            </div>
            <div class="request-sub" style="color: #475569; margin-top: 4px;">
              <ng-container *ngIf="showWorkerName">
                {{ shift.start_time | date:'MMM d, y, h:mm a' }} → {{ shift.end_time | date:'h:mm a' }}
              </ng-container>
              <ng-container *ngIf="!showWorkerName">
                Status: {{ shift.status }}
              </ng-container>
            </div>
          </div>
          
          <div style="display: flex; gap: 8px;">
            
            <button *ngIf="showDelete" class="btn ghost dark" style="padding: 6px 12px; font-size: 0.85rem; border: 1px solid #dc2626; color: #dc2626; cursor: pointer;" 
                    (click)="onDelete.emit(shift.id || shift.shift_id)">
              <mat-icon style="font-size: 16px; width: 16px; height: 16px; vertical-align: middle;">delete</mat-icon> Delete
            </button>

            <ng-container *ngIf="showActions">
              <button class="btn ghost dark" style="padding: 6px 12px; font-size: 0.85rem; border: 1px solid #cbd5e1;" 
                      (click)="onReject.emit(shift.shift_id)">
                Reject
              </button>
              <button class="btn primary" style="padding: 6px 12px; font-size: 0.85rem;" 
                      (click)="onApprove.emit(shift)">
                Approve
              </button>
            </ng-container>
            
          </div>

        </div>
      </li>
    </ul>
  `
})
export class ShiftListComponent {
  @Input() shifts: any[] = [];
  @Input() isLoading: boolean = false;
  @Input() loadingMessage: string = 'Loading...';
  @Input() emptyMessage: string = 'No shifts found.';
  @Input() showWorkerName: boolean = false;
  @Input() showActions: boolean = false;
  @Input() showDelete: boolean = false;

  @Output() onApprove = new EventEmitter<any>();
  @Output() onReject = new EventEmitter<number>();
  @Output() onDelete = new EventEmitter<number>();
}