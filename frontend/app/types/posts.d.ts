interface PostMeta {
  id: number;
  cover: string;
  title: string;
  summary: string;
  publishedAt: string;
  updatedAt: string;
  readTimeMinutes: number;
  viewCount: number;
  tags: [];
}

interface Post extends PostMeta {
  content: string;
}
