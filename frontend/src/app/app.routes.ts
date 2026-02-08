import { Routes } from '@angular/router';
import { LandingComponent } from './landing/landing';
import { WorkerDashboardComponent } from './worker-dashboard/worker-dashboard';

export const routes: Routes = [
    { path:'login', component: LandingComponent },
    { path: 'worker-dashboard', component: WorkerDashboardComponent }
];
