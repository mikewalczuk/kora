import { createFileRoute } from "@tanstack/react-router";
import { NoteDetail } from "@/features/notes/NoteDetail";

export const Route = createFileRoute("/notes/$noteId")({
  component: NoteDetailPage,
});

function NoteDetailPage() {
  const { noteId } = Route.useParams();
  return (
    <div className="mx-auto mt-16 max-w-lg px-4">
      <NoteDetail noteId={noteId} />
    </div>
  );
}
