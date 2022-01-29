import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'kp-textile';
  constructor(){
    localStorage.setItem('token', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicmFuY2hJZCI6ImJyYW5jaElkIiwidXNlcm5hbWUiOiJVc2VybmFtZSIsInJvbGVzIjpbMV0sImV4cCI6MTY0ODU4NTk4NSwiaWF0IjoxNjQzNDg4Mzg1fQ.k_GUKEKbdYKiHki_1BZe77E7Or6zmSZkd3YTnobpOPU')
  }
}
