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
  HStack,
  Button,
  IconButton,
  Select,
  Flex,
  Spacer,
} from "@chakra-ui/react";
import { ChevronLeftIcon, ChevronRightIcon } from "@chakra-ui/icons";
import type { CountryRevenue } from "../services/api";

interface CountryRevenueTableProps {
  data: CountryRevenue[];
  totalCount: number;
  currentPage: number;
  pageSize: number;
  onPageChange: (page: number, pageSize: number) => void;
  isLoading?: boolean;
}

export function CountryRevenueTable({
  data,
  totalCount,
  currentPage,
  pageSize,
  onPageChange,
  isLoading = false,
}: CountryRevenueTableProps) {
  const totalPages = Math.ceil(totalCount / pageSize);
  const startIndex = (currentPage - 1) * pageSize + 1;
  const endIndex = Math.min(currentPage * pageSize, totalCount);

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
      onPageChange(newPage, pageSize);
    }
  };

  const handlePageSizeChange = (newPageSize: number) => {
    if (newPageSize !== pageSize) {
      onPageChange(1, newPageSize); // Always go to page 1 when changing page size
    }
  };

  const generatePageNumbers = () => {
    const pages = [];
    const maxVisiblePages = 5;

    if (totalPages <= maxVisiblePages) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      if (currentPage <= 3) {
        pages.push(1, 2, 3, 4, "...", totalPages);
      } else if (currentPage >= totalPages - 2) {
        pages.push(
          1,
          "...",
          totalPages - 3,
          totalPages - 2,
          totalPages - 1,
          totalPages
        );
      } else {
        pages.push(
          1,
          "...",
          currentPage - 1,
          currentPage,
          currentPage + 1,
          "...",
          totalPages
        );
      }
    }

    return pages;
  };

  return (
    <Box
      bg="white"
      p={6}
      borderRadius="lg"
      shadow="md"
      border="1px"
      borderColor="gray.200"
    >
      <Flex mb={4} align="center">
        <Box>
          <Heading size="lg" mb={2}>
            Country-Level Revenue
          </Heading>
          <Text color="gray.600" fontSize="sm">
            Showing {startIndex} to {endIndex} of {totalCount.toLocaleString()}{" "}
            country-product combinations
          </Text>
        </Box>
        <Spacer />
        <HStack spacing={2}>
          <Text fontSize="sm" color="gray.600">
            Show:
          </Text>
          <Select
            size="sm"
            width="auto"
            value={pageSize}
            onChange={(e) => handlePageSizeChange(Number(e.target.value))}
            isDisabled={isLoading}
          >
            <option value={25}>25</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
            <option value={200}>200</option>
          </Select>
          <Text fontSize="sm" color="gray.600">
            per page
          </Text>
        </HStack>
      </Flex>

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
            {data.map((item, index) => (
              <Tr
                key={`${item.country}-${item.product_name}-${index}`}
                _hover={{ bg: "gray.50" }}
                opacity={isLoading ? 0.6 : 1}
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

      {/* Pagination Controls */}
      <Flex mt={4} justify="space-between" align="center">
        <HStack spacing={2}>
          <IconButton
            aria-label="Previous page"
            icon={<ChevronLeftIcon />}
            size="sm"
            isDisabled={currentPage === 1 || isLoading}
            onClick={() => handlePageChange(currentPage - 1)}
          />

          {generatePageNumbers().map((page, index) =>
            page === "..." ? (
              <Text key={index} px={2} color="gray.500">
                ...
              </Text>
            ) : (
              <Button
                key={index}
                size="sm"
                variant={currentPage === page ? "solid" : "ghost"}
                colorScheme={currentPage === page ? "blue" : "gray"}
                isDisabled={isLoading}
                onClick={() => handlePageChange(page as number)}
                minW="40px"
              >
                {page}
              </Button>
            )
          )}

          <IconButton
            aria-label="Next page"
            icon={<ChevronRightIcon />}
            size="sm"
            isDisabled={currentPage === totalPages || isLoading}
            onClick={() => handlePageChange(currentPage + 1)}
          />
        </HStack>

        <Text fontSize="sm" color="gray.600">
          Page {currentPage} of {totalPages}
        </Text>
      </Flex>
    </Box>
  );
}
