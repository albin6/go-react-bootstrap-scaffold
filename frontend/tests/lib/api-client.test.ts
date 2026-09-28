import { describe, it, expect } from 'vitest';
import { ApiError } from '@/lib/api-client';

describe('ApiClient & ApiError', () => {
  it('creates ApiError with standard code, message, and status', () => {
    const err = new ApiError('NOT_FOUND', 'Entity not found', 404, [{ field: 'id', issue: 'missing' }]);
    expect(err.name).toBe('ApiError');
    expect(err.code).toBe('NOT_FOUND');
    expect(err.message).toBe('Entity not found');
    expect(err.status).toBe(404);
    expect(err.details).toHaveLength(1);
  });
});
