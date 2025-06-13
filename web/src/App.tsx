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
import { useState, useCallback } from "react";
import {
  getAnalytics,
  getStats,
  getCountryRevenue,
  type AnalyticsResponse,
  type StatsResponse,
  type CountryRevenuePaginatedResponse,
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
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(50);

  const {
    data: analytics,
    isLoading: analyticsLoading,
    error: analyticsError,
    refetch: refetchAnalytics,
  } = useQuery<AnalyticsResponse>({
    queryKey: ["analytics"],
    queryFn: getAnalytics,
  });

  const { data: stats } = useQuery<StatsResponse>({
    queryKey: ["analytics-stats"],
    queryFn: getStats,
  });

  const {
    data: countryRevenueData,
    isLoading: countryRevenueLoading,
    error: countryRevenueError,
    refetch: refetchCountryRevenue,
  } = useQuery<CountryRevenuePaginatedResponse, Error>({
    queryKey: ["country-revenue", currentPage, pageSize],
    queryFn: () => getCountryRevenue(pageSize, (currentPage - 1) * pageSize),
    staleTime: 30000,
  });

  const handlePageChange = useCallback(
    (page: number, newPageSize: number) => {
      if (newPageSize !== pageSize) {
        setPageSize(newPageSize);
        setCurrentPage(1);
      } else {
        setCurrentPage(page);
      }
    },
    [pageSize]
  );

  const handleRefresh = useCallback(() => {
    refetchAnalytics();
    refetchCountryRevenue();
  }, [refetchAnalytics, refetchCountryRevenue]);

  const isLoading = analyticsLoading && countryRevenueLoading;
  const hasError = analyticsError || countryRevenueError;

  return (
    <Container maxW="container.xl" py={8}>
      <VStack spacing={8} align="stretch">
        <Box>
          <Heading size="xl" mb={4}>
            ABT Analytics Dashboard
          </Heading>
          <RefreshButton onRefresh={handleRefresh} />
        </Box>

        <DebugInfo />

        {isLoading && <LoadingSpinner />}

        {hasError && (
          <Alert status="error">
            <AlertIcon />
            <VStack align="start" spacing={2}>
              <Text>Failed to load analytics data. Please try again.</Text>
              <Code fontSize="sm">
                Error:{" "}
                {analyticsError instanceof Error
                  ? analyticsError.message
                  : countryRevenueError instanceof Error
                  ? countryRevenueError.message
                  : "Unknown error"}
              </Code>
              <Text fontSize="sm" color="gray.600">
                Make sure your Go backend is running on localhost:8080
              </Text>
            </VStack>
          </Alert>
        )}

        {/* Stats Card */}
        {stats && <StatsCard stats={stats} />}

        {/* Country Revenue Table with Pagination */}
        {countryRevenueData && (
          <CountryRevenueTable
            data={countryRevenueData.data}
            totalCount={countryRevenueData.total}
            currentPage={currentPage}
            pageSize={pageSize}
            onPageChange={handlePageChange}
            isLoading={countryRevenueLoading}
          />
        )}

        {/* Charts using analytics data */}
        {analytics && (
          <>
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
