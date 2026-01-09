import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { BlogResponse } from '../../models/models';

@Component({
  selector: 'app-blog-detail',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './blog-detail.component.html',
  styleUrl: './blog-detail.component.css'
})
export class BlogDetailComponent implements OnInit {
  blog: BlogResponse | null = null;
  loading = true;
  error: string | null = null;

  constructor(
    private apiService: ApiService,
    private route: ActivatedRoute
  ) { }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      const id = params['id'];
      this.loadBlog(id);
    });
  }

  loadBlog(id: number): void {
    this.apiService.getBlogById(id).subscribe({
      next: (data) => {
        this.blog = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading blog:', err);
        this.error = 'Failed to load blog post';
        this.loading = false;
      }
    });
  }
}
