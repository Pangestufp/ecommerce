import { NotebookPen } from "lucide-react";

/**
 * NoteSection
 * Catatan dari pembeli ke penjual.
 *
 * Props:
 *  - value      : string, isi catatan
 *  - onChange   : (value: string) => void
 */
export default function NoteSection({ value, onChange }) {
  return (
    <div className="bg-white rounded-2xl border border-gray-100 px-4 py-4 mb-3">
      <div className="flex items-center gap-2 mb-3">
        <NotebookPen size={14} className="text-blue-500 shrink-0" />
        <h3 className="text-sm font-semibold text-gray-800">Catatan untuk Penjual</h3>
      </div>

      <textarea
        rows={3}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Tulis catatan atau permintaan khusus (opsional)"
        maxLength={300}
        className="w-full text-sm text-gray-800 placeholder-gray-300 bg-gray-50 border border-gray-200 rounded-xl px-3 py-2.5 resize-none focus:outline-none focus:ring-2 focus:ring-blue-100 focus:border-blue-400 transition"
      />

      <p className="text-[11px] text-gray-300 text-right mt-1">
        {value.length}/300
      </p>
    </div>
  );
}