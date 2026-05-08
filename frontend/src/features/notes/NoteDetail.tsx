import { Link, useRouter } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { useNote, useDeleteNote, useUpdateNote } from "./api";
import { useCreatePractice } from "@/features/practices/api";

interface Props {
  noteId: string;
}

export function NoteDetail({ noteId }: Props) {
  const { note, isLoading, isError } = useNote(noteId);
  const { deleteNote, isPending: isDeleting } = useDeleteNote();
  const { updateNote, isPending: isSaving } = useUpdateNote();
  const { createPractice, isPending: isPracticing } = useCreatePractice();
  const router = useRouter();

  async function handlePractice() {
    const practice = await createPractice(noteId);
    router.navigate({
      to: "/practice/$practiceId",
      params: { practiceId: practice.id },
      state: { practice } as never,
    });
  }

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  useEffect(() => {
    if (note) {
      setTitle(note.title);
      setContent(note.content);
    }
  }, [note]);

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === "s") {
        e.preventDefault();
        updateNote(noteId, title, content);
      }
    }
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [noteId, title, content]);

  if (isLoading) return <p className="text-sm text-gray-400">Loading…</p>;
  if (isError || !note) return <p className="text-sm text-red-500">Note not found.</p>;

  async function handleDelete() {
    await deleteNote(noteId);
    router.navigate({ to: "/" });
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <Link to="/" className="text-sm text-gray-500 hover:text-gray-800">
          ← Back
        </Link>
        <div className="flex gap-3">
          <button
            onClick={() => updateNote(noteId, title, content)}
            disabled={isSaving}
            className="text-sm text-violet-600 hover:text-violet-800 disabled:opacity-50"
          >
            {isSaving ? "Saving…" : "Save"}
          </button>
          <button
            onClick={handlePractice}
            disabled={isPracticing}
            className="text-sm text-violet-600 hover:text-violet-800 disabled:opacity-50"
          >
            {isPracticing ? "Creating…" : "Practice"}
          </button>
          <button
            onClick={handleDelete}
            disabled={isDeleting}
            className="text-sm text-gray-400 hover:text-red-500 disabled:opacity-50"
          >
            {isDeleting ? "Deleting…" : "Delete"}
          </button>
        </div>
      </div>
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className="text-xl font-semibold text-gray-900 outline-none"
        placeholder="Title"
      />
      <p className="text-xs text-gray-400">
        {new Date(note.created_at).toLocaleDateString()}
      </p>
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={12}
        className="w-full resize-none text-sm text-gray-700 outline-none"
        placeholder="Write something…"
      />
    </div>
  );
}
