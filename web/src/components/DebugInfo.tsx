import {
  Box,
  Button,
  Text,
  VStack,
  Alert,
  AlertIcon,
  Code,
} from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import { healthCheck } from "../services/api";

interface HealthResponse {
  status: string;
}

export function DebugInfo() {
  const {
    data: health,
    isLoading: healthLoading,
    error: healthError,
    refetch: refetchHealth,
  } = useQuery<HealthResponse>({
    queryKey: ["health"],
    queryFn: healthCheck,
    retry: false,
  });

  return (
    <Box bg="gray.50" p={4} borderRadius="md" mb={4}>
      <Text fontWeight="bold" mb={2}>
        Debug Information
      </Text>
      <VStack align="start" spacing={2}>
        <Text fontSize="sm">
          <strong>API Base URL:</strong>{" "}
          <Code>
            {import.meta.env.DEV
              ? "/api/v1 (proxied)"
              : "http://localhost:8080/api/v1"}
          </Code>
        </Text>

        <Box>
          <Text fontSize="sm" mb={1}>
            <strong>Backend Health Check:</strong>
          </Text>
          {healthLoading && (
            <Text fontSize="sm" color="blue.500">
              Checking...
            </Text>
          )}
          {healthError && (
            <Alert status="error" size="sm">
              <AlertIcon />
              <Text fontSize="sm">
                Backend not accessible:{" "}
                {healthError instanceof Error
                  ? healthError.message
                  : "Unknown error"}
              </Text>
            </Alert>
          )}
          {health && (
            <Alert status="success" size="sm">
              <AlertIcon />
              <Text fontSize="sm">Backend is healthy: {health.status}</Text>
            </Alert>
          )}
          <Button size="xs" mt={1} onClick={() => refetchHealth()}>
            Test Connection
          </Button>
        </Box>
      </VStack>
    </Box>
  );
}
