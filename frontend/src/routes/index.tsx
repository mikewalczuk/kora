import { createFileRoute, useRouter } from "@tanstack/react-router";
import { Button } from "@/components/ui/Button";
import { useLogout } from "@/features/auth/api";
import { useCreateNote } from "@/features/notes/api";
import { NotesList } from "@/features/notes/NotesList";

export const Route = createFileRoute("/")({
  component: HomePage,
});

function HomePage() {
  const router = useRouter();
  const logout = useLogout();
  const { createNote, isPending: isCreating } = useCreateNote();

  function handleLogout() {
    logout.mutate(undefined, {
      onSuccess: () => router.navigate({ to: "/auth/login" }),
    });
  }

  async function handleNewNote() {
    const note = await createNote("Untitled", "");
    router.navigate({ to: "/notes/$noteId", params: { noteId: note.id } });
  }

  return (
    <div className="mx-auto mt-16 flex max-w-lg flex-col gap-6 px-4">
      <div className="flex items-center justify-between">
        <h1 className="text-lg font-semibold">Notes</h1>
        <div className="flex gap-2">
          <Button onClick={handleNewNote} disabled={isCreating}>
            {isCreating ? "Creating…" : "New note"}
          </Button>
          <Button variant="ghost" onClick={handleLogout} disabled={logout.isPending}>
            {logout.isPending ? "Logging out…" : "Log out"}
          </Button>
        </div>
      </div>
      <NotesList />
    </div>
  );
}
