import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from '../../services/api.service';
import { HomeResponse } from '../../models/models';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent implements OnInit {
  homeData: HomeResponse | null = null;
  loading = true;
  error: string | null = null;

  constructor(private apiService: ApiService) { }

  ngOnInit(): void {
    this.loadHomeData();
  }

  loadHomeData(): void {
    this.apiService.getHome().subscribe({
      next: (data) => {
        this.homeData = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading home data:', err);
        this.error = 'Failed to load home data';
        this.loading = false;
      }
    });
  }
}
