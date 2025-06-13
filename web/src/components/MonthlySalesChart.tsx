import { Box, Heading, Text } from "@chakra-ui/react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import type { MonthlySales } from "../services/api";

interface MonthlySalesChartProps {
  data: MonthlySales[];
}

export function MonthlySalesChart({ data }: MonthlySalesChartProps) {
  const chartData = data.map((item) => ({
    month: item.month,
    sales_volume: item.sales_volume,
    item_count: item.item_count,
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
        Monthly Sales Volume
      </Heading>
      <Text color="gray.600" mb={6}>
        Sales trends over time
      </Text>

      <Box h="400px">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart
            data={chartData}
            margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="month" />
            <YAxis yAxisId="left" />
            <YAxis yAxisId="right" orientation="right" />
            <Tooltip
              formatter={(value: number, name: string) => [
                name === "sales_volume"
                  ? `$${value.toLocaleString(undefined, {
                      minimumFractionDigits: 2,
                    })}`
                  : value.toLocaleString(),
                name === "sales_volume" ? "Sales Volume" : "Items Sold",
              ]}
            />
            <Legend />
            <Line
              yAxisId="left"
              type="monotone"
              dataKey="sales_volume"
              stroke="#3182ce"
              strokeWidth={3}
              name="Sales Volume ($)"
            />
            <Line
              yAxisId="right"
              type="monotone"
              dataKey="item_count"
              stroke="#38a169"
              strokeWidth={2}
              name="Items Sold"
            />
          </LineChart>
        </ResponsiveContainer>
      </Box>
    </Box>
  );
}
