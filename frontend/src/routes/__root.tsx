import { createRootRoute, Outlet, redirect } from '@tanstack/react-router';
import { getMe } from '@/features/auth/api';

export const Route = createRootRoute({
  beforeLoad: async ({ location }) => {
    if (location.pathname.startsWith('/auth')) return;
    try {
      await getMe();
    } catch {
      throw redirect({ to: '/auth/login' });
    }
  },
  component: () => (
    <div className="min-h-svh">
      <Outlet />
    </div>
  ),
});
