"use client";

import { useState, useEffect } from 'react';
import { apiService } from '@/lib/api';

export function useStreak() {
  const [streak, setStreak] = useState<number>(0);
  const [lastEntryDate, setLastEntryDate] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStreak = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await apiService.getStreak();
      setStreak(data.streak);
      setLastEntryDate(data.last_entry_date || null);
    } catch (err: any) {
      setError(err.message || 'Failed to load streak');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStreak();
  }, []);

  return { streak, lastEntryDate, loading, error, refetch: fetchStreak };
}
