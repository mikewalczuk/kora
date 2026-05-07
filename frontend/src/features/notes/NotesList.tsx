import { Link } from "@tanstack/react-router";
import { useNotesList, useDeleteNote } from "./api";

export function NotesList() {
  const { notes, isLoading, isError } = useNotesList();
  const { deleteNote, isPending: isDeleting } = useDeleteNote();

  if (isLoading) return <p className="text-sm text-gray-400">Loading notes…</p>;
  if (isError) return <p className="text-sm text-red-500">Failed to load notes.</p>;
  if (notes.length === 0) return <p className="text-sm text-gray-400">No notes yet.</p>;

  return (
    <ul className="flex flex-col gap-2">
      {notes.map((note) => (
        <li key={note.id} className="flex items-center rounded-md border border-gray-200 hover:bg-gray-50">
          <Link
            to="/notes/$noteId"
            params={{ noteId: note.id }}
            className="flex-1 p-3"
          >
            <p className="text-sm font-medium text-gray-900">{note.title}</p>
            <p className="mt-1 text-xs text-gray-400">
              {new Date(note.created_at).toLocaleDateString()}
            </p>
          </Link>
          <button
            onClick={() => deleteNote(note.id)}
            disabled={isDeleting}
            className="px-3 py-3 text-gray-400 hover:text-red-500 disabled:opacity-50"
            aria-label="Delete note"
          >
            ✕
          </button>
        </li>
      ))}
    </ul>
  );
}
