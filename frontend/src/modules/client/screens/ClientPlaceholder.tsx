interface ClientPlaceholderProps {
  title: string;
}

export function ClientPlaceholder({ title }: ClientPlaceholderProps) {
  return (
    <div className="p-4 flex items-center justify-center h-64 text-gray-400 text-[13px]">
      {title} — coming soon
    </div>
  );
}
