import { createFileRoute, Link } from "@tanstack/react-router";
import type { Exercise, Practice } from "@/lib/api/generated/koraAPI.schemas";
import { usePractice } from "@/features/practices/api";

export const Route = createFileRoute("/practice/$practiceId")({
  component: PracticePage,
});

function PracticePage() {
  const { practiceId } = Route.useParams();
  const { practice, isLoading, isError } = usePractice(practiceId);

  return (
    <div className="mx-auto mt-16 max-w-lg px-4 flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <Link to="/" className="text-sm text-gray-500 hover:text-gray-800">
          ← Back
        </Link>
        <span className="text-xs text-gray-400">Practice</span>
      </div>

      {isLoading && <p className="text-sm text-gray-400">Loading…</p>}
      {isError && <p className="text-sm text-red-500">Practice not found.</p>}
      {practice && <PracticeDetail practice={practice} />}
    </div>
  );
}

function PracticeDetail({ practice }: { practice: Practice }) {
  return (
    <div className="flex flex-col gap-6">
      {practice.exercises.map((ex, i) => (
        <ExerciseCard key={ex.id} exercise={ex} index={i} />
      ))}
    </div>
  );
}

function ExerciseCard({ exercise, index }: { exercise: Exercise; index: number }) {
  return (
    <div className="flex flex-col gap-4">
      <p className="text-xs text-gray-400 uppercase tracking-wide">
        Exercise {index + 1}
      </p>
      {exercise.questions.map((q) => (
        <div key={q.id} className="flex flex-col gap-3">
          <p className="text-sm font-medium text-gray-800">{q.question}</p>
          <ul className="flex flex-col gap-2">
            {q.options.map((opt) => (
              <li key={opt.id}>
                <button className="w-full text-left rounded border border-gray-200 px-3 py-2 text-sm text-gray-700 hover:border-violet-400 hover:bg-violet-50 transition-colors">
                  {opt.text}
                </button>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
}
