interface LoadingSkeletonProps {
  rows?: number;
  cols: number;
}

export function LoadingSkeleton({ rows = 5, cols }: LoadingSkeletonProps) {
  return (
    <>
      {Array.from({ length: rows }).map((_, i) => (
        <tr key={i} className="border-b border-gray-100">
          {Array.from({ length: cols }).map((_, j) => (
            <td key={j} className="p-4">
              <div className="h-3 bg-gray-100 rounded animate-pulse" style={{ width: `${60 + ((i + j) % 3) * 15}%` }} />
            </td>
          ))}
        </tr>
      ))}
    </>
  );
}
