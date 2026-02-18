export interface ApiResponse<T> {
  code: number;
  msg: string;
  data: T;
}

export interface ApiListResponse<T> extends ApiResponse<T[]> {}

export interface ApiPageData<T> {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ApiPageResponse<T> extends ApiResponse<ApiPageData<T>> {}

export interface ApiError {
  code: number;
  msg: string;
}
