package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joaovictorsl/dcache/client"
)

const (
	CMD_SET    = "SET"
	CMD_GET    = "GET"
	CMD_HAS    = "HAS"
	CMD_DELETE = "DELETE"
)

func main() {
	nodeaddr := flag.String("nodeaddr", ":3000", "node address")
	flag.Parse()

	c := client.New(*nodeaddr)
	err := c.Connect(5, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer c.End()

	startCliClient(c)
}

func startCliClient(c *client.DCacheClient) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		cmdStr := scanner.Text()

		if strings.ToUpper(cmdStr) == "EXIT" {
			break
		}

		res, err := handleCmd(cmdStr, c)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(res)
	}
}

func handleCmd(cmdStr string, c *client.DCacheClient) (string, error) {
	if len(cmdStr) == 0 {
		return "", fmt.Errorf("empty command")
	}

	parts := strings.Split(cmdStr, " ")

	switch strings.ToUpper(parts[0]) {
	case CMD_SET:
		return handleSet(parts, c)

	case CMD_GET:
		return handleGet(parts, c)

	case CMD_HAS:
		return handleHas(parts, c)

	case CMD_DELETE:
		return handleDelete(parts, c)

	default:
		return "", fmt.Errorf("invalid command")
	}
}

func handleSet(parts []string, c *client.DCacheClient) (string, error) {
	if len(parts) != 4 {
		return "", fmt.Errorf("INVALID SET COMMAND")
	}

	key := parts[1]
	value := parts[2]
	ttlStr := parts[3]
	ttl64, err := strconv.ParseUint(ttlStr, 10, 32)
	if err != nil {
		return "", err
	}

	if err := c.Set(key, []byte(value), uint32(ttl64)); err != nil {
		return "", err
	}

	return fmt.Sprintf("SET COMMAND KEY (%s) OK", key), nil
}

func handleGet(parts []string, c *client.DCacheClient) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("INVALID GET COMMAND")
	}

	key := parts[1]

	res, found, err := c.Get(key)
	if err != nil {
		return "", err
	}

	if !found {
		return fmt.Sprintf("GET COMMAND KEY (%s) NOT FOUND", key), nil
	}

	return fmt.Sprintf("GET COMMAND KEY (%s) FOUND\n%s", key, string(res)), nil
}

func handleDelete(parts []string, c *client.DCacheClient) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("INVALID DELETE COMMAND")
	}

	key := parts[1]

	if err := c.Delete(key); err != nil {
		return "", err
	}

	return fmt.Sprintf("DELETE (%s) COMMAND OK", key), nil
}

func handleHas(parts []string, c *client.DCacheClient) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("INVALID HAS COMMAND")
	}

	key := parts[1]

	ok, err := c.Has(key)
	if err != nil {
		return "", err
	}

	if ok {
		return fmt.Sprintf("HAS COMMAND KEY (%s) FOUND", key), nil
	}

	return fmt.Sprintf("HAS COMMAND KEY (%s) NOT FOUND", key), nil
}
