// Standardized API Envelope Definitions

export interface ApiResponse<T> {
  success: true;
  data: T;
}

export interface ApiErrorDetail {
  field?: string;
  issue: string;
}

export interface ApiErrorInfo {
  code: string;
  message: string;
  details?: ApiErrorDetail[];
}

export interface ApiErrorResponse {
  success: false;
  error: ApiErrorInfo;
}

export type ApiResult<T> = ApiResponse<T> | ApiErrorResponse;

export interface PaginationMetadata {
  page: number;
  limit: number;
  total_items: number;
  total_pages: number;
  has_next: boolean;
}

export interface PaginatedData<T> {
  items: T[];
  pagination: PaginationMetadata;
}
