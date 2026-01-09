export interface HomeResponse {
  title: string;
  description: string;
  name: string;
}

export interface ProfileResponse {
  name: string;
  title: string;
  bio: string;
  email: string;
  skills: string[];
  github: string;
  linkedin: string;
}

export interface BlogResponse {
  id: number;
  title: string;
  content: string;
  author: string;
  date: string;
  tags: string[];
}

export interface BlogRequest {
  title: string;
  content: string;
  author: string;
  tags: string[];
}
