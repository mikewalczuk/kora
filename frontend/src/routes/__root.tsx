import { createRootRoute, Outlet, redirect } from '@tanstack/react-router';
import { getMe } from '@/features/auth/api';
import { SSEProvider } from '@/lib/sse/SSEProvider';
import { Toaster } from '@/lib/sse/Toaster';

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
    <SSEProvider>
      <div className="min-h-svh">
        <Outlet />
      </div>
      <Toaster />
    </SSEProvider>
  ),
});
