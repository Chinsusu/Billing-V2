export function trimStringFilters<T extends Record<string, string>>(filters: T): T {
  return Object.fromEntries(
    Object.entries(filters).map(([key, value]) => [key, value.trim()]),
  ) as T;
}

export function hasActiveFilters(filters: Record<string, string>): boolean {
  return Object.values(filters).some((value) => value.trim() !== "");
}

export function includesFilter(value: string | number | null | undefined, query: string): boolean {
  const trimmed = query.trim().toLowerCase();
  if (!trimmed) {
    return true;
  }
  return String(value ?? "").toLowerCase().includes(trimmed);
}

export function equalsFilter(value: string | number | null | undefined, query: string): boolean {
  const trimmed = query.trim().toLowerCase();
  if (!trimmed) {
    return true;
  }
  return String(value ?? "").toLowerCase() === trimmed;
}

export function matchesAmountRange(value: number, minValue: string, maxValue: string): boolean {
  const min = minValue.trim() === "" ? null : Number(minValue);
  const max = maxValue.trim() === "" ? null : Number(maxValue);

  if (min !== null && !Number.isNaN(min) && value < min) {
    return false;
  }
  if (max !== null && !Number.isNaN(max) && value > max) {
    return false;
  }
  return true;
}
