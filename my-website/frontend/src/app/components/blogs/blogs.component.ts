import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { BlogResponse } from '../../models/models';

@Component({
  selector: 'app-blogs',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './blogs.component.html',
  styleUrl: './blogs.component.css'
})
export class BlogsComponent implements OnInit {
  blogs: BlogResponse[] = [];
  loading = true;
  error: string | null = null;

  constructor(private apiService: ApiService) { }

  ngOnInit(): void {
    this.loadBlogs();
  }

  loadBlogs(): void {
    this.apiService.getBlogs().subscribe({
      next: (data) => {
        this.blogs = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading blogs:', err);
        this.error = 'Failed to load blogs';
        this.loading = false;
      }
    });
  }
}
