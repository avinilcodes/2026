import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { HomeResponse, ProfileResponse, BlogResponse, BlogRequest } from '../models/models';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) { }

  getHome(): Observable<HomeResponse> {
    return this.http.get<HomeResponse>(`${this.apiUrl}/home`);
  }

  getProfile(): Observable<ProfileResponse> {
    return this.http.get<ProfileResponse>(`${this.apiUrl}/profile`);
  }

  getBlogs(): Observable<BlogResponse[]> {
    return this.http.get<BlogResponse[]>(`${this.apiUrl}/blogs`);
  }

  getBlogById(id: number): Observable<BlogResponse> {
    return this.http.get<BlogResponse>(`${this.apiUrl}/blogs/${id}`);
  }

  createBlog(blog: BlogRequest): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/blogs`, blog);
  }
}
