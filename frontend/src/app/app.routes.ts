import { Routes } from '@angular/router';
import { LandingComponent } from './landing/landing.component';
import { WorkerDashboardComponent } from './worker-dashboard/worker-dashboard.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { SupervisorDashboardComponent } from './supervisor-dashboard/supervisor-dashboard.component';
import { authGuard } from './services/auth.guard';

export const routes: Routes = [
    { path:'', redirectTo:'/home', pathMatch: 'full' },
    { path:'home', component: LandingComponent },
    { path:'login', component:LoginComponent },
    { path:'register', component:RegisterComponent },
    { path: 'worker-dashboard', component: WorkerDashboardComponent, canActivate: [authGuard] },
    { path: 'supervisor-dashboard', component: SupervisorDashboardComponent, canActivate: [authGuard] }
];
