import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from '@/components/ui/Button';

describe('Button Component', () => {
  it('renders button with children text', () => {
    render(<Button>Click Me</Button>);
    expect(screen.getByRole('button', { name: /click me/i })).toBeInTheDocument();
  });

  it('handles click events', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Trigger</Button>);
    fireEvent.click(screen.getByRole('button', { name: /trigger/i }));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('disables button when disabled prop is passed', () => {
    render(<Button disabled>Disabled Action</Button>);
    expect(screen.getByRole('button', { name: /disabled action/i })).toBeDisabled();
  });

  it('renders spinner and disables button when isLoading is true', () => {
    render(<Button isLoading>Processing</Button>);
    const btn = screen.getByRole('button', { name: /processing/i });
    expect(btn).toBeDisabled();
    expect(btn.querySelector('svg')).toBeInTheDocument();
  });
});
