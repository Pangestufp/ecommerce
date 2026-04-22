import Button from "../ui/Button";

export default function Table({ columns, data, rowKey, onUpdate, onDelete }) {
  const alignMap = {
    left: "text-left",
    center: "text-center",
    right: "text-right",
  };

  if (!data || data.length === 0) {
    return (
      <div className="bg-white border border-gray-200 rounded-xl p-8 text-center text-sm text-gray-400">
        No data available
      </div>
    );
  }

  return (
    <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              {columns.map((col) => (
                <th
                  key={col.key}
                  className={`px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide ${
                    alignMap[col.align || "left"]
                  }`}
                >
                  {col.label}
                </th>
              ))}

              {(onUpdate || onDelete) && (
                <th className="px-4 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wide">
                  Action
                </th>
              )}
            </tr>
          </thead>

          <tbody className="divide-y divide-gray-100">
            {data.map((row) => (
              <tr
                key={row[rowKey]}
                className="hover:bg-gray-50 transition-colors"
              >
                {columns.map((col) => (
                  <td
                    key={col.key}
                    className={`px-4 py-3 text-gray-700 ${
                      alignMap[col.align || "left"]
                    }`}
                  >
                    {col.render
                      ? col.render(row[col.key], row)
                      : row[col.key]}
                  </td>
                ))}

                {(onUpdate || onDelete) && (
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-2">
                      {onUpdate && (
                        <Button
                          variant="secondary"
                          onClick={() => onUpdate(row)}
                          className="text-xs px-3 py-1"
                        >
                          Edit
                        </Button>
                      )}

                      {onDelete && (
                        <Button
                          variant="danger"
                          onClick={() => onDelete(row[rowKey])}
                          className="text-xs px-3 py-1"
                        >
                          Hapus
                        </Button>
                      )}
                    </div>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}