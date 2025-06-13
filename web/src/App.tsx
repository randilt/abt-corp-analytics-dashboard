import {
  Box,
  Container,
  Heading,
  VStack,
  Alert,
  AlertIcon,
  Code,
  Text,
} from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import {
  getAnalytics,
  getStats,
  type AnalyticsResponse,
  type StatsResponse,
} from "./services/api";
import { LoadingSpinner } from "./components/LoadingSpinner";
import { StatsCard } from "./components/StatsCard";
import { CountryRevenueTable } from "./components/CountryRevenueTable";
import { TopProductsChart } from "./components/TopProductsChart";
import { MonthlySalesChart } from "./components/MonthlySalesChart";
import { TopRegionsChart } from "./components/TopRegionsChart";
import { RefreshButton } from "./components/RefreshButton";
import { DebugInfo } from "./components/DebugInfo";

function App() {
  const {
    data: analytics,
    isLoading,
    error,
    refetch,
  } = useQuery<AnalyticsResponse>({
    queryKey: ["analytics"],
    queryFn: getAnalytics,
  });

  const { data: stats } = useQuery<StatsResponse>({
    queryKey: ["analytics-stats"],
    queryFn: getStats,
  });

  return (
    <Container maxW="container.xl" py={8}>
      <VStack spacing={8} align="stretch">
        <Box>
          <Heading size="xl" mb={4}>
            ABT Analytics Dashboard
          </Heading>
          <RefreshButton onRefresh={() => refetch()} />
        </Box>

        <DebugInfo />

        {isLoading && <LoadingSpinner />}

        {error && (
          <Alert status="error">
            <AlertIcon />
            <VStack align="start" spacing={2}>
              <Text>Failed to load analytics data. Please try again.</Text>
              <Code fontSize="sm">
                Error:{" "}
                {error instanceof Error ? error.message : "Unknown error"}
              </Code>
              <Text fontSize="sm" color="gray.600">
                Make sure your Go backend is running on localhost:8080
              </Text>
            </VStack>
          </Alert>
        )}

        {analytics && (
          <>
            {stats && <StatsCard stats={stats} />}
            <CountryRevenueTable data={analytics.country_revenue} />
            <TopProductsChart data={analytics.top_products} />
            <MonthlySalesChart data={analytics.monthly_sales} />
            <TopRegionsChart data={analytics.top_regions} />
          </>
        )}
      </VStack>
    </Container>
  );
}

export default App;
