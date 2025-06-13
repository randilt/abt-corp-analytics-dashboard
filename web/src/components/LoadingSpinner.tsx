import { Box, Spinner, Text, VStack } from "@chakra-ui/react";

export function LoadingSpinner() {
  return (
    <Box
      minH="100vh"
      display="flex"
      alignItems="center"
      justifyContent="center"
    >
      <VStack spacing={4}>
        <Spinner size="xl" color="blue.500" thickness="4px" />
        <Text color="gray.600">Loading analytics data...</Text>
      </VStack>
    </Box>
  );
}
