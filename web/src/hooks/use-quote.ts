"use client";

import { useState, useEffect } from 'react';
import { apiService } from '@/lib/api';
import { Quote } from '@/lib/types';

export function useRandomQuote() {
  const [quote, setQuote] = useState<Quote | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRandomQuote = async () => {
    setLoading(true);
    setError(null);

    try {
      const newQuote = await apiService.getRandomQuote();
      setQuote(newQuote);
    } catch (err: any) {
      setError(err.message || 'Failed to load quote');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRandomQuote();
  }, []);

  return {
    quote,
    loading,
    error,
    refetch: fetchRandomQuote,
  };
}
