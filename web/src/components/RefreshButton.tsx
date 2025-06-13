import { Button, useToast } from "@chakra-ui/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { refreshCache } from "../services/api";

interface RefreshButtonProps {
  onRefresh?: () => void;
}

export function RefreshButton({ onRefresh }: RefreshButtonProps) {
  const toast = useToast();
  const queryClient = useQueryClient();

  const { mutate: refresh, isPending } = useMutation({
    mutationFn: refreshCache,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["analytics"] });
      queryClient.invalidateQueries({ queryKey: ["analytics-stats"] });
      toast({
        title: "Cache refreshed",
        status: "success",
        duration: 2000,
        isClosable: true,
      });
      onRefresh?.();
    },
    onError: (error) => {
      toast({
        title: "Failed to refresh cache",
        description: error instanceof Error ? error.message : "Unknown error",
        status: "error",
        duration: 5000,
        isClosable: true,
      });
    },
  });

  return (
    <Button
      colorScheme="blue"
      size="sm"
      onClick={() => refresh()}
      isLoading={isPending}
      loadingText="Refreshing..."
    >
      Refresh Data
    </Button>
  );
}
