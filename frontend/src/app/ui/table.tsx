"use client";
import { useTheme } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Table from '@mui/material/Table';
import TableHead from '@mui/material/TableHead';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableFooter from '@mui/material/TableFooter';
import TablePagination from '@mui/material/TablePagination';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import IconButton from '@mui/material/IconButton';
import FirstPageIcon from '@mui/icons-material/FirstPage';
import KeyboardArrowLeft from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRight from '@mui/icons-material/KeyboardArrowRight';
import LastPageIcon from '@mui/icons-material/LastPage';
import { useEffect, useRef, useState } from "react";

interface TablePaginationActionsProps {
    count: number;
    page: number;
    rowsPerPage: number;
    onPageChange: (
        event: React.MouseEvent<HTMLButtonElement>,
        newPage: number,
    ) => void;
}

function TablePaginationActions(props: TablePaginationActionsProps) {
    const theme = useTheme();
    const { count, page, rowsPerPage, onPageChange } = props;

    const handleFirstPageButtonClick = (
        event: React.MouseEvent<HTMLButtonElement>,
    ) => {
        onPageChange(event, 0);
    };

    const handleBackButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        onPageChange(event, page - 1);
    };

    const handleNextButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        onPageChange(event, page + 1);
    };

    const handleLastPageButtonClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        onPageChange(event, Math.max(0, Math.ceil(count / rowsPerPage) - 1));
    };

    return (
        <Box sx={{ flexShrink: 0, ml: 2.5 }}>
            <IconButton
                onClick={handleFirstPageButtonClick}
                disabled={page === 0}
                aria-label="first page"
            >
                {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
            </IconButton>
            <IconButton
                onClick={handleBackButtonClick}
                disabled={page === 0}
                aria-label="previous page"
            >
                {theme.direction === 'rtl' ? <KeyboardArrowRight /> : <KeyboardArrowLeft />}
            </IconButton>
            <IconButton
                onClick={handleNextButtonClick}
                disabled={page >= Math.ceil(count / rowsPerPage) - 1}
                aria-label="next page"
            >
                {theme.direction === 'rtl' ? <KeyboardArrowLeft /> : <KeyboardArrowRight />}
            </IconButton>
            <IconButton
                onClick={handleLastPageButtonClick}
                disabled={page >= Math.ceil(count / rowsPerPage) - 1}
                aria-label="last page"
            >
                {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
            </IconButton>
        </Box>
    );
}

type Word = {
    id: number;
    word: string;
    definition: string;
}

const API_ENDPOINT = "http://localhost:8080";

async function getWords(limit: number, page: number): Promise<Word[]> {
    console.log("calllllllll api");
    // const res = await fetch(API_ENDPOINT + "/words?limit=" + limit + "?page=" + page);
    // const data = res.json();
    // return data;
    return [
        { id: 1, word: "hello", definition: "is hello" },
        { id: 2, word: "hello2", definition: "is hello2" },
        { id: 3, word: "hello3", definition: "is hello3" },
        { id: 4, word: "hello4", definition: "is hello4" },
        { id: 5, word: "hello4", definition: "is hello4" },
        { id: 6, word: "hello4", definition: "is hello4" },
        { id: 7, word: "hello4", definition: "is hello4" },
        { id: 8, word: "hello4", definition: "is hello4" },
        { id: 9, word: "hello4", definition: "is hello4" },
        { id: 10, word: "hello4", definition: "is hello4" },
    ];
}


export default function MyTable() {
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(5);
    const [words, setWords] = useState<Word[]>([]);

    // for page load
    useEffect(() => {
        console.log("use effect 1");
        const fetchWords = async () => {
            const result = await getWords(rowsPerPage, page);
            setWords(result);
        };
        fetchWords();
    }, []);

    // Avoid a layout jump when reaching the last page with empty rows.
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - words.length) : 0;

    const handleChangePage = async (
        event: React.MouseEvent<HTMLButtonElement> | null,
        newPage: number,
    ) => {
        const words = await getWords(rowsPerPage, newPage);
        setPage(newPage);
        setWords(words);
    };

    const handleChangeRowsPerPage = async (
        event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        const newRowsPerPage = parseInt(event.target.value, 10);
        const newPage = 0;
        const words = await getWords(newRowsPerPage, newPage);
        setRowsPerPage(newRowsPerPage);
        setPage(newPage);
        setWords(words);
    };

    return (
        <TableContainer component={Paper} style={{ width: "70%" }}>
            <Table>
                <TableHead>
                    <TableRow style={{ backgroundColor: "black" }}>
                        <TableCell style={{ color: "white" }} align="left">Word</TableCell>
                        <TableCell style={{ color: "white" }} align="left">Definition</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {words.map((word) => (
                        // two columns
                        <TableRow key={word.id}>
                            <TableCell component="th" style={{ width: 100 }}>
                                {word.word}
                            </TableCell>
                            <TableCell>
                                {word.definition}
                            </TableCell>
                        </TableRow>
                    ))}
                    {/* for last page */}
                    {emptyRows > 0 && (
                        <TableRow style={{ height: 53 * emptyRows }}>
                            <TableCell colSpan={6} />
                        </TableRow>
                    )}
                </TableBody>
                <TableFooter>
                    <TableRow>
                        <TablePagination
                            rowsPerPageOptions={[5, 10, 25, { label: 'All', value: -1 }]}
                            colSpan={3}
                            count={words.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            onPageChange={handleChangePage}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                            ActionsComponent={TablePaginationActions}
                        />
                    </TableRow>
                </TableFooter>
            </Table>
        </TableContainer>
    );
}