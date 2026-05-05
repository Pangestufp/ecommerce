export default function Section({ title, action, children }) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div className="flex items-center justify-between px-5 py-3.5 border-b border-gray-100">
        <h2 className="text-sm font-semibold text-gray-700">{title}</h2>
        {action && <div className="flex items-center gap-3">{action}</div>}
      </div>
      <div className="p-5">{children}</div>
    </div>
  );
}