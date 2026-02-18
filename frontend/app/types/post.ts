export interface PostMeta {
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

export interface Post extends PostMeta {
  content: string;
}

export interface PostInput {
  cover: string;
  summary: string;
  title: string;
  content: string;
  tagIds: number[];
  status: "draft" | "published";
}
