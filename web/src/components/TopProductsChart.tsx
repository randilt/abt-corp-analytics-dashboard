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
import type { ProductFrequency } from "../services/api";

interface TopProductsChartProps {
  data: ProductFrequency[];
}

export function TopProductsChart({ data }: TopProductsChartProps) {
  const chartData = data.map((item) => ({
    name:
      item.product_name.length > 20
        ? item.product_name.substring(0, 20) + "..."
        : item.product_name,
    fullName: item.product_name,
    purchase_count: item.purchase_count,
    current_stock: item.current_stock,
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
        Top 20 Most Purchased Products
      </Heading>
      <Text color="gray.600" mb={6}>
        Purchase count vs. current stock levels
      </Text>

      <Box h="500px">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={chartData}
            margin={{ top: 20, right: 30, left: 20, bottom: 100 }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="name"
              angle={-45}
              textAnchor="end"
              height={120}
              interval={0}
              fontSize={12}
            />
            <YAxis />
            <Tooltip<number, "Purchase Count" | "Current Stock">
              labelFormatter={(label, payload) =>
                payload?.[0]?.payload?.fullName || label
              }
              formatter={(
                value: number,
                name: "Purchase Count" | "Current Stock"
              ) => [value.toLocaleString(), name]}
            />
            <Legend />
            <Bar
              dataKey="purchase_count"
              fill="#3182ce"
              name="Purchase Count"
            />
            <Bar dataKey="current_stock" fill="#38a169" name="Current Stock" />
          </BarChart>
        </ResponsiveContainer>
      </Box>
    </Box>
  );
}
