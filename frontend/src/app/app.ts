import { Component} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { LandingComponent } from './landing/landing';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, MatButtonModule, LandingComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class AppComponent {
  title = 'frontend';
}

export { AppComponent as App};
