export default function Table({ columns, data, rowKey, onUpdate, onDelete }) {
  if (!data || data.length === 0) {
    return (
      <div className="text-center py-8 text-sm text-gray-400">
        No data available
      </div>
    );
  }

  return (
    <div className="w-full overflow-x-auto">
      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b border-gray-200">
            {columns.map((col) => (
              <th key={col.key} className="text-left px-3 py-2.5 text-xs font-medium text-gray-400 tracking-wide">
                {col.label}
              </th>
            ))}
            <th />
          </tr>
        </thead>
        <tbody>
          {data.map((row) => (
            <tr key={row[rowKey]} className="border-b border-gray-100 hover:bg-gray-50">
              {columns.map((col) => (
                <td key={col.key} className="px-3 py-2.5 text-gray-800">
                  {col.render ? col.render(row[col.key], row) : row[col.key]}
                </td>
              ))}
              <td className="px-3 py-2.5">
                <div className="flex gap-2">
                  {onUpdate && (
                    <button
                      onClick={() => onUpdate(row)}
                      className="text-xs font-medium px-2.5 py-1 rounded border border-gray-200 text-gray-500 hover:bg-gray-100 hover:text-gray-800 transition-colors"
                    >
                      Edit
                    </button>
                  )}
                  {onDelete && (
                    <button
                      onClick={() => onDelete(row[rowKey])}
                      className="text-xs font-medium px-2.5 py-1 rounded border border-gray-200 text-gray-500 hover:bg-red-50 hover:text-red-700 hover:border-red-200 transition-colors"
                    >
                      Delete
                    </button>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}