import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Heading,
  Text,
  Badge,
} from "@chakra-ui/react";
import type { CountryRevenue } from "../services/api";

interface CountryRevenueTableProps {
  data: CountryRevenue[];
}

export function CountryRevenueTable({ data }: CountryRevenueTableProps) {
  const displayData = data.slice(0, 50); // Show first 50 for better performance

  return (
    <Box
      bg="white"
      p={6}
      borderRadius="lg"
      shadow="md"
      border="1px"
      borderColor="gray.200"
    >
      <Heading size="lg" mb={4}>
        Country-Level Revenue
      </Heading>
      <Text color="gray.600" mb={4}>
        Showing top {displayData.length} of {data.length} country-product
        combinations
      </Text>

      <TableContainer maxH="600px" overflowY="auto">
        <Table variant="simple" size="sm">
          <Thead bg="gray.50" position="sticky" top={0} zIndex={1}>
            <Tr>
              <Th>Country</Th>
              <Th>Product Name</Th>
              <Th isNumeric>Total Revenue</Th>
              <Th isNumeric>Transactions</Th>
            </Tr>
          </Thead>
          <Tbody>
            {displayData.map((item, index) => (
              <Tr
                key={`${item.country}-${item.product_name}-${index}`}
                _hover={{ bg: "gray.50" }}
              >
                <Td>
                  <Badge colorScheme="blue" variant="subtle">
                    {item.country}
                  </Badge>
                </Td>
                <Td>
                  <Text fontSize="sm" noOfLines={2}>
                    {item.product_name}
                  </Text>
                </Td>
                <Td isNumeric fontWeight="medium">
                  $
                  {item.total_revenue.toLocaleString(undefined, {
                    minimumFractionDigits: 2,
                  })}
                </Td>
                <Td isNumeric>{item.transaction_count.toLocaleString()}</Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </TableContainer>
    </Box>
  );
}
