import { Inbox } from "lucide-react";

interface EmptyStateProps {
  title: string;
  description?: string;
  action?: React.ReactNode;
}

export function EmptyState({ title, description, action }: EmptyStateProps) {
  return (
    <tr>
      <td colSpan={99} className="py-16 text-center">
        <div className="flex flex-col items-center gap-2">
          <Inbox size={28} className="text-gray-300" />
          <span className="text-[13px] font-medium text-gray-500">{title}</span>
          {description && <span className="text-[12px] text-gray-400">{description}</span>}
          {action && <div className="mt-2">{action}</div>}
        </div>
      </td>
    </tr>
  );
}
