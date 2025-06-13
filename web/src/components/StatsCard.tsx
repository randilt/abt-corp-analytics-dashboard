import {
  Box,
  Grid,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Badge,
} from "@chakra-ui/react";
import type { StatsResponse } from "../services/api";

interface StatsCardProps {
  stats: StatsResponse;
}

export function StatsCard({ stats }: StatsCardProps) {
  return (
    <Box
      bg="white"
      p={6}
      borderRadius="lg"
      shadow="md"
      border="1px"
      borderColor="gray.200"
    >
      <Grid templateColumns="repeat(auto-fit, minmax(200px, 1fr))" gap={6}>
        <Stat>
          <StatLabel>Total Records</StatLabel>
          <StatNumber>{stats.total_records.toLocaleString()}</StatNumber>
          <StatHelpText>
            <Badge colorScheme={stats.cache_hit ? "green" : "blue"}>
              {stats.cache_hit ? "Cache Hit" : "Fresh Load"}
            </Badge>
          </StatHelpText>
        </Stat>

        <Stat>
          <StatLabel>Processing Time</StatLabel>
          <StatNumber>{stats.processing_time_ms}ms</StatNumber>
          <StatHelpText>Data processing duration</StatHelpText>
        </Stat>

        <Stat>
          <StatLabel>Countries</StatLabel>
          <StatNumber>{stats.country_revenue_count}</StatNumber>
          <StatHelpText>Unique country-product pairs</StatHelpText>
        </Stat>

        <Stat>
          <StatLabel>Top Products</StatLabel>
          <StatNumber>{stats.top_products_count}</StatNumber>
          <StatHelpText>Most purchased items</StatHelpText>
        </Stat>
      </Grid>
    </Box>
  );
}
