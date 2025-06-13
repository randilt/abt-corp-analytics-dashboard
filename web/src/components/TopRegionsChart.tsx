import { Box, Heading, Text } from "@chakra-ui/react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import type { RegionRevenue } from "../services/api";

interface TopRegionsChartProps {
  data: RegionRevenue[];
}

export function TopRegionsChart({ data }: TopRegionsChartProps) {
  const chartData = data.slice(0, 30).map((item) => ({
    region:
      item.region.length > 15
        ? item.region.substring(0, 15) + "..."
        : item.region,
    fullRegion: item.region,
    total_revenue: item.total_revenue,
    items_sold: item.items_sold,
  }));

  return (
    <Box
      bg="white"
      p={6}
      borderRadius="lg"
      shadow="md"
      border="1px"
      borderColor="gray.200"
    >
      <Heading size="lg" mb={2}>
        Top Regions by Revenue
      </Heading>
      <Text color="gray.600" mb={6}>
        Revenue and items sold by region (showing top 30, Available{" "}
        {data.length})
      </Text>

      <Box h="500px">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={chartData}
            margin={{ top: 20, right: 30, left: 20, bottom: 100 }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="region"
              angle={-45}
              textAnchor="end"
              height={120}
              interval={0}
              fontSize={12}
            />
            <YAxis yAxisId="left" />
            <YAxis yAxisId="right" orientation="right" />
            <Tooltip<number, "Total Revenue" | "Items Sold">
              labelFormatter={(label, payload) =>
                payload?.[0]?.payload?.fullRegion || label
              }
              formatter={(
                value: number,
                name: "Total Revenue" | "Items Sold"
              ) => [
                name === "Total Revenue"
                  ? `$${value.toLocaleString(undefined, {
                      minimumFractionDigits: 2,
                    })}`
                  : value.toLocaleString(),
                name,
              ]}
            />
            <Legend />
            <Bar
              yAxisId="left"
              dataKey="total_revenue"
              fill="#3182ce"
              name="Total Revenue ($)"
            />
            <Bar
              yAxisId="right"
              dataKey="items_sold"
              fill="#38a169"
              name="Items Sold"
            />
          </BarChart>
        </ResponsiveContainer>
      </Box>
    </Box>
  );
}
